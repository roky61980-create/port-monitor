package store

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/roky61980-create/port-monitor/internal/model"
	_ "modernc.org/sqlite"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil { return nil, err }
	db.SetMaxOpenConns(1)
	if err := db.Ping(); err != nil { db.Close(); return nil, err }
	s := &Store{db: db}
	if err := s.migrate(); err != nil { db.Close(); return nil, err }
	return s, nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) migrate() error {
	const schema = `
	CREATE TABLE IF NOT EXISTS devices (id TEXT PRIMARY KEY,name TEXT,ip TEXT UNIQUE,vendor TEXT,model TEXT,location TEXT,description TEXT,snmp_version TEXT,enabled INTEGER,status TEXT,last_seen DATETIME);
	CREATE TABLE IF NOT EXISTS ports (id TEXT PRIMARY KEY,device_id TEXT,if_index INTEGER,if_name TEXT,if_descr TEXT,if_alias TEXT,admin_status TEXT,oper_status TEXT,speed INTEGER,duplex TEXT,last_change DATETIME,alert_status TEXT,UNIQUE(device_id,if_index));
	CREATE TABLE IF NOT EXISTS events (id TEXT PRIMARY KEY,timestamp DATETIME,device_id TEXT,device TEXT,vendor TEXT,model TEXT,interface TEXT,interface_id INTEGER,event_type TEXT,severity TEXT,oid TEXT,value TEXT,message TEXT,raw_data TEXT);
	CREATE TABLE IF NOT EXISTS alerts (id TEXT PRIMARY KEY,device_id TEXT,device TEXT,interface TEXT,alert_type TEXT,severity TEXT,status TEXT,current_value TEXT,threshold TEXT,acknowledged_by TEXT,created_at DATETIME,updated_at DATETIME,resolved_at DATETIME);
	CREATE TABLE IF NOT EXISTS optical (id TEXT PRIMARY KEY,device_id TEXT,device TEXT,interface TEXT,module TEXT,rx_power REAL,rx_status TEXT,tx_power REAL,tx_status TEXT,temperature REAL,voltage REAL,bias REAL,vendor TEXT,part_number TEXT,serial_number TEXT,wavelength TEXT,status TEXT,last_update DATETIME);
	CREATE INDEX IF NOT EXISTS idx_events_time ON events(timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_events_device ON events(device_id,timestamp DESC);
	CREATE INDEX IF NOT EXISTS idx_alerts_status ON alerts(status,updated_at DESC);`
	_, err := s.db.Exec(schema)
	return err
}

func (s *Store) UpsertDevice(d model.Device) error {
	_, err := s.db.Exec(`INSERT INTO devices(id,name,ip,vendor,model,location,description,snmp_version,enabled,status,last_seen) VALUES(?,?,?,?,?,?,?,?,?,?,?)
	ON CONFLICT(id) DO UPDATE SET name=excluded.name,ip=excluded.ip,vendor=excluded.vendor,model=excluded.model,location=excluded.location,description=excluded.description,snmp_version=excluded.snmp_version,enabled=excluded.enabled,status=excluded.status,last_seen=excluded.last_seen`, d.ID,d.Name,d.IP,d.Vendor,d.Model,d.Location,d.Description,d.SNMPVersion,boolInt(d.Enabled),d.Status,d.LastSeen)
	return err
}

func boolInt(v bool) int { if v { return 1 }; return 0 }

func (s *Store) Devices() ([]model.Device,error) {
	rows,err:=s.db.Query(`SELECT id,name,ip,vendor,model,location,description,snmp_version,enabled,status,last_seen FROM devices ORDER BY name`); if err!=nil{return nil,err}; defer rows.Close()
	var out []model.Device
	for rows.Next(){var d model.Device;var enabled int;var seen sql.NullTime;if err:=rows.Scan(&d.ID,&d.Name,&d.IP,&d.Vendor,&d.Model,&d.Location,&d.Description,&d.SNMPVersion,&enabled,&d.Status,&seen);err!=nil{return nil,err};d.Enabled=enabled==1;if seen.Valid{d.LastSeen=&seen.Time};out=append(out,d)}
	return out,rows.Err()
}

func (s *Store) Device(id string)(model.Device,error){var d model.Device;var enabled int;var seen sql.NullTime;err:=s.db.QueryRow(`SELECT id,name,ip,vendor,model,location,description,snmp_version,enabled,status,last_seen FROM devices WHERE id=?`,id).Scan(&d.ID,&d.Name,&d.IP,&d.Vendor,&d.Model,&d.Location,&d.Description,&d.SNMPVersion,&enabled,&d.Status,&seen);if err!=nil{return d,err};d.Enabled=enabled==1;if seen.Valid{d.LastSeen=&seen.Time};return d,nil}

func (s *Store) DeleteDevice(id string) error { _,err:=s.db.Exec(`DELETE FROM devices WHERE id=?`,id);return err }

func (s *Store) SaveEvent(e model.NetworkEvent) error {value,_:=json.Marshal(e.Value);raw,_:=json.Marshal(e.RawData);_,err:=s.db.Exec(`INSERT OR IGNORE INTO events(id,timestamp,device_id,device,vendor,model,interface,interface_id,event_type,severity,oid,value,message,raw_data) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Timestamp,e.DeviceID,e.Device,e.Vendor,e.Model,e.Interface,e.InterfaceID,e.EventType,e.Severity,e.OID,string(value),e.Message,string(raw));return err}

func (s *Store) Events(limit,offset int)([]model.NetworkEvent,error){rows,err:=s.db.Query(`SELECT id,timestamp,device_id,device,vendor,model,interface,interface_id,event_type,severity,oid,value,message,raw_data FROM events ORDER BY timestamp DESC LIMIT ? OFFSET ?`,limit,offset);if err!=nil{return nil,err};defer rows.Close();var out []model.NetworkEvent;for rows.Next(){var e model.NetworkEvent;var value,raw string;if err:=rows.Scan(&e.ID,&e.Timestamp,&e.DeviceID,&e.Device,&e.Vendor,&e.Model,&e.Interface,&e.InterfaceID,&e.EventType,&e.Severity,&e.OID,&value,&e.Message,&raw);err!=nil{return nil,err};json.Unmarshal([]byte(value),&e.Value);json.Unmarshal([]byte(raw),&e.RawData);out=append(out,e)};return out,rows.Err()}

func (s *Store) AlertByKey(deviceID,iface,alertType string)(model.Alert,error){var a model.Alert;var resolved sql.NullTime;err:=s.db.QueryRow(`SELECT id,device_id,device,interface,alert_type,severity,status,current_value,threshold,acknowledged_by,created_at,updated_at,resolved_at FROM alerts WHERE device_id=? AND interface=? AND alert_type=? AND status!='RESOLVED' ORDER BY created_at DESC LIMIT 1`,deviceID,iface,alertType).Scan(&a.ID,&a.DeviceID,&a.Device,&a.Interface,&a.AlertType,&a.Severity,&a.Status,&a.CurrentValue,&a.Threshold,&a.AcknowledgedBy,&a.CreatedAt,&a.UpdatedAt,&resolved);if err!=nil{return a,err};if resolved.Valid{a.ResolvedAt=&resolved.Time};return a,nil}

func (s *Store) SaveAlert(a model.Alert)error{_,err:=s.db.Exec(`INSERT INTO alerts(id,device_id,device,interface,alert_type,severity,status,current_value,threshold,acknowledged_by,created_at,updated_at,resolved_at) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET severity=excluded.severity,status=excluded.status,current_value=excluded.current_value,threshold=excluded.threshold,acknowledged_by=excluded.acknowledged_by,updated_at=excluded.updated_at,resolved_at=excluded.resolved_at`,a.ID,a.DeviceID,a.Device,a.Interface,a.AlertType,a.Severity,a.Status,a.CurrentValue,a.Threshold,a.AcknowledgedBy,a.CreatedAt,a.UpdatedAt,a.ResolvedAt);return err}

func (s *Store) Alerts()( []model.Alert,error){rows,err:=s.db.Query(`SELECT id,device_id,device,interface,alert_type,severity,status,current_value,threshold,COALESCE(acknowledged_by,''),created_at,updated_at,resolved_at FROM alerts ORDER BY updated_at DESC LIMIT 500`);if err!=nil{return nil,err};defer rows.Close();var out []model.Alert;for rows.Next(){var a model.Alert;var resolved sql.NullTime;if err:=rows.Scan(&a.ID,&a.DeviceID,&a.Device,&a.Interface,&a.AlertType,&a.Severity,&a.Status,&a.CurrentValue,&a.Threshold,&a.AcknowledgedBy,&a.CreatedAt,&a.UpdatedAt,&resolved);err!=nil{return nil,err};if resolved.Valid{a.ResolvedAt=&resolved.Time};out=append(out,a)};return out,rows.Err()}

func (s *Store) Alert(id string)(model.Alert,error){var a model.Alert;var resolved sql.NullTime;err:=s.db.QueryRow(`SELECT id,device_id,device,interface,alert_type,severity,status,current_value,threshold,COALESCE(acknowledged_by,''),created_at,updated_at,resolved_at FROM alerts WHERE id=?`,id).Scan(&a.ID,&a.DeviceID,&a.Device,&a.Interface,&a.AlertType,&a.Severity,&a.Status,&a.CurrentValue,&a.Threshold,&a.AcknowledgedBy,&a.CreatedAt,&a.UpdatedAt,&resolved);if err!=nil{return a,err};if resolved.Valid{a.ResolvedAt=&resolved.Time};if errors.Is(err,sql.ErrNoRows){return a,sql.ErrNoRows};return a,nil}

func (s *Store) SavePort(p model.Port) error {_,err:=s.db.Exec(`INSERT INTO ports(id,device_id,if_index,if_name,if_descr,if_alias,admin_status,oper_status,speed,duplex,last_change,alert_status) VALUES(?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET if_name=excluded.if_name,if_descr=excluded.if_descr,if_alias=excluded.if_alias,admin_status=excluded.admin_status,oper_status=excluded.oper_status,speed=excluded.speed,duplex=excluded.duplex,last_change=excluded.last_change,alert_status=excluded.alert_status`,p.ID,p.DeviceID,p.IfIndex,p.IfName,p.IfDescr,p.IfAlias,p.AdminStatus,p.OperStatus,p.Speed,p.Duplex,p.LastChange,p.AlertStatus);return err}

func (s *Store) Ports()([]model.Port,error){rows,err:=s.db.Query(`SELECT id,device_id,if_index,if_name,if_descr,if_alias,admin_status,oper_status,speed,duplex,last_change,alert_status FROM ports ORDER BY device_id,if_index`);if err!=nil{return nil,err};defer rows.Close();var out []model.Port;for rows.Next(){var p model.Port;var lc sql.NullTime;if err:=rows.Scan(&p.ID,&p.DeviceID,&p.IfIndex,&p.IfName,&p.IfDescr,&p.IfAlias,&p.AdminStatus,&p.OperStatus,&p.Speed,&p.Duplex,&lc,&p.AlertStatus);err!=nil{return nil,err};if lc.Valid{p.LastChange=&lc.Time};out=append(out,p)};return out,rows.Err()}

func (s *Store) SaveOptical(r model.OpticalReading) error {_,err:=s.db.Exec(`INSERT INTO optical(id,device_id,device,interface,module,rx_power,rx_status,tx_power,tx_status,temperature,voltage,bias,vendor,part_number,serial_number,wavelength,status,last_update) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?) ON CONFLICT(id) DO UPDATE SET rx_power=excluded.rx_power,rx_status=excluded.rx_status,tx_power=excluded.tx_power,tx_status=excluded.tx_status,temperature=excluded.temperature,voltage=excluded.voltage,bias=excluded.bias,vendor=excluded.vendor,part_number=excluded.part_number,serial_number=excluded.serial_number,wavelength=excluded.wavelength,status=excluded.status,last_update=excluded.last_update`,r.ID,r.DeviceID,r.Device,r.Interface,r.Module,r.RXPower,r.RXStatus,r.TXPower,r.TXStatus,r.Temperature,r.Voltage,r.Bias,r.Vendor,r.PartNumber,r.SerialNumber,r.Wavelength,r.Status,r.LastUpdate);return err}

func (s *Store) Optical()([]model.OpticalReading,error){rows,err:=s.db.Query(`SELECT id,device_id,device,interface,module,rx_power,rx_status,tx_power,tx_status,temperature,voltage,bias,vendor,part_number,serial_number,wavelength,status,last_update FROM optical ORDER BY device,interface`);if err!=nil{return nil,err};defer rows.Close();var out []model.OpticalReading;for rows.Next(){var r model.OpticalReading;if err:=rows.Scan(&r.ID,&r.DeviceID,&r.Device,&r.Interface,&r.Module,nullableFloat(&r.RXPower),&r.RXStatus,nullableFloat(&r.TXPower),&r.TXStatus,nullableFloat(&r.Temperature),nullableFloat(&r.Voltage),nullableFloat(&r.Bias),&r.Vendor,&r.PartNumber,&r.SerialNumber,&r.Wavelength,&r.Status,&r.LastUpdate);err!=nil{return nil,err};out=append(out,r)};return out,rows.Err()}

func nullableFloat(dst **float64) any { return &nullableFloatScanner{dst:dst} }
type nullableFloatScanner struct{dst **float64}
func(n *nullableFloatScanner)Scan(v any)error{if v==nil{*n.dst=nil;return nil};var f float64;switch x:=v.(type){case float64:f=x;case int64:f=float64(x);default:return errors.New("invalid float")};*n.dst=&f;return nil}

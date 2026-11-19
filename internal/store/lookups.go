package store

import (
	"database/sql"
	"errors"
	"github.com/roky61980-create/port-monitor/internal/model"
)

var ErrNotFound = sql.ErrNoRows

func (s *Store) DeviceByIP(ip string)(model.Device,error){var d model.Device;var enabled int;err:=s.db.QueryRow(`SELECT id,name,ip,vendor,model,location,description,snmp_version,enabled,status,last_seen FROM devices WHERE ip=?`,ip).Scan(&d.ID,&d.Name,&d.IP,&d.Vendor,&d.Model,&d.Location,&d.Description,&d.SNMPVersion,&enabled,&d.Status,&d.LastSeen);if errors.Is(err,sql.ErrNoRows){return d,ErrNotFound};d.Enabled=enabled==1;return d,err}
func (s *Store) PortByIndex(deviceID string,index int)(model.Port,error){var p model.Port;err:=s.db.QueryRow(`SELECT id,device_id,if_index,if_name,if_descr,if_alias,admin_status,oper_status,speed,duplex,last_change,alert_status FROM ports WHERE device_id=? AND if_index=?`,deviceID,index).Scan(&p.ID,&p.DeviceID,&p.IfIndex,&p.IfName,&p.IfDescr,&p.IfAlias,&p.AdminStatus,&p.OperStatus,&p.Speed,&p.Duplex,&p.LastChange,&p.AlertStatus);return p,err}

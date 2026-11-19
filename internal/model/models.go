package model

import "time"

type Device struct {
	ID string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	IP string `json:"ip" db:"ip"`
	Vendor string `json:"vendor" db:"vendor"`
	Model string `json:"model" db:"model"`
	Location string `json:"location" db:"location"`
	Description string `json:"description" db:"description"`
	SNMPVersion string `json:"snmp_version" db:"snmp_version"`
	Enabled bool `json:"enabled" db:"enabled"`
	Status string `json:"status" db:"status"`
	LastSeen *time.Time `json:"last_seen,omitempty" db:"last_seen"`
}

type Port struct {
	ID string `json:"id" db:"id"`
	DeviceID string `json:"device_id" db:"device_id"`
	IfIndex int `json:"if_index" db:"if_index"`
	IfName string `json:"if_name" db:"if_name"`
	IfDescr string `json:"if_descr" db:"if_descr"`
	IfAlias string `json:"if_alias" db:"if_alias"`
	AdminStatus string `json:"admin_status" db:"admin_status"`
	OperStatus string `json:"oper_status" db:"oper_status"`
	Speed int64 `json:"speed" db:"speed"`
	Duplex string `json:"duplex" db:"duplex"`
	LastChange *time.Time `json:"last_change,omitempty" db:"last_change"`
	AlertStatus string `json:"alert_status" db:"alert_status"`
}

type NetworkEvent struct {
	ID string `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	DeviceID string `json:"device_id"`
	Device string `json:"device"`
	Vendor string `json:"vendor"`
	Model string `json:"model"`
	Interface string `json:"interface"`
	InterfaceID int `json:"interface_id"`
	EventType string `json:"event_type"`
	Severity string `json:"severity"`
	OID string `json:"oid"`
	Value any `json:"value,omitempty"`
	Message string `json:"message"`
	RawData any `json:"raw_data,omitempty"`
}

type Alert struct {
	ID string `json:"id"`
	DeviceID string `json:"device_id"`
	Device string `json:"device"`
	Interface string `json:"interface"`
	AlertType string `json:"alert_type"`
	Severity string `json:"severity"`
	Status string `json:"status"`
	CurrentValue string `json:"current_value"`
	Threshold string `json:"threshold"`
	AcknowledgedBy string `json:"acknowledged_by,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

type OpticalReading struct {
	ID string `json:"id"`
	DeviceID string `json:"device_id"`
	Device string `json:"device"`
	Interface string `json:"interface"`
	Module string `json:"module"`
	RXPower *float64 `json:"rx_power,omitempty"`
	RXStatus string `json:"rx_status"`
	TXPower *float64 `json:"tx_power,omitempty"`
	TXStatus string `json:"tx_status"`
	Temperature *float64 `json:"temperature,omitempty"`
	Voltage *float64 `json:"voltage,omitempty"`
	Bias *float64 `json:"bias,omitempty"`
	Vendor string `json:"vendor"`
	PartNumber string `json:"part_number,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
	Wavelength string `json:"wavelength,omitempty"`
	Status string `json:"status"`
	LastUpdate time.Time `json:"last_update"`
}

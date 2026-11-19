package store

import (
	"database/sql"
	"github.com/roky61980-create/port-monitor/internal/model"
)
func (s *Store) AlertByKeySafe(deviceID,iface,alertType string)(model.Alert,error){a,e:=s.AlertByKey(deviceID,iface,alertType);if e!=nil && e!=sql.ErrNoRows{return a,e};return a,e}

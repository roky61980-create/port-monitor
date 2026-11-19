package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const DefaultPath = "config.yaml"

type Config struct {
	Server struct {
		Listen string `yaml:"listen"`
	} `yaml:"server"`
	SNMP struct {
		TrapListen string `yaml:"trap_listen"`
	} `yaml:"snmp"`
	Optical struct {
		Enabled      bool          `yaml:"enabled"`
		PollInterval time.Duration `yaml:"poll_interval"`
		Workers      int           `yaml:"workers"`
	} `yaml:"optical"`
	OTel struct {
		Endpoint string `yaml:"endpoint"`
	} `yaml:"otel"`
	InfluxDB struct {
		URL    string `yaml:"url"`
		Org    string `yaml:"org"`
		Bucket string `yaml:"bucket"`
		Token  string `yaml:"token"`
	} `yaml:"influxdb"`
	Alerts struct {
		DedupWindow time.Duration `yaml:"dedup_window"`
		Hysteresis  bool          `yaml:"hysteresis"`
	} `yaml:"alerts"`
	SSE struct{ Enabled bool `yaml:"enabled"` } `yaml:"sse"`
	Database struct{ Path string `yaml:"path"` } `yaml:"database"`
}

func Load(path string) (Config, error) {
	var c Config
	if path == "" { path = DefaultPath }
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) { return c, err }
	if err == nil {
		if err := yaml.Unmarshal(data, &c); err != nil { return c, err }
	}
	applyDefaults(&c)
	if v := os.Getenv("INFLUXDB_TOKEN"); v != "" { c.InfluxDB.Token = v }
	return c, nil
}

func applyDefaults(c *Config) {
	if c.Server.Listen == "" { c.Server.Listen = "0.0.0.0:8080" }
	if c.SNMP.TrapListen == "" { c.SNMP.TrapListen = "0.0.0.0:162" }
	if c.Optical.PollInterval == 0 { c.Optical.PollInterval = time.Minute }
	if c.Optical.Workers == 0 { c.Optical.Workers = 4 }
	if c.OTel.Endpoint == "" { c.OTel.Endpoint = "127.0.0.1:4317" }
	if c.Alerts.DedupWindow == 0 { c.Alerts.DedupWindow = time.Minute }
	if c.Database.Path == "" { c.Database.Path = "port-monitor.db" }
	if c.InfluxDB.URL == "" { c.InfluxDB.URL = "http://127.0.0.1:8086" }
	if c.InfluxDB.Org == "" { c.InfluxDB.Org = "emnetwork" }
	if c.InfluxDB.Bucket == "" { c.InfluxDB.Bucket = "network" }
}

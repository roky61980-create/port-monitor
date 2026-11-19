# EM Network Port Monitor

Focused network-device port/link alerts, SNMP trap ingestion, optical DDM provider framework, SSE dashboard, SQLite state and OTLP trace/event export. No RMON or traffic accounting is included.

## Run

Install Go 1.24+, copy `config.yaml`, set secrets in the process environment, then `go run ./cmd/server`. Web/API defaults to TCP/8080; traps default to UDP/162 (requires suitable capabilities or privileged service configuration). SQLite stores application records. The browser UI is served from `web/`.

## Scope and support

SNMPv2c notifications are accepted. Standard IF-MIB linkUp/linkDown notifications are normalized; device identity is resolved from configured device IPs. Provider interfaces exist for Cisco, MikroTik, BDCOM and generic devices. Optical provider OIDs are intentionally unconfigured: no fabricated measurements are produced. See [docs/OPTICAL-DDM.md](docs/OPTICAL-DDM.md).

OTLP traces/events are exported to the configured Collector endpoint. `deploy/otel/otel-collector.yaml` is a minimal local debug pipeline. Vector/InfluxDB are deployment stubs and require adapting transforms and line-protocol mapping before production use.

## Documentation

- [Architecture](ARCHITECTURE.md)
- [Configuration](CONFIGURATION.md)
- [Deployment](DEPLOYMENT.md)
- SNMP [traps](docs/SNMP-TRAP.md), [ports](docs/PORT-ALERT.md), [optical](docs/OPTICAL-DDM.md), [OTel](docs/OPENTELEMETRY.md), [Vector](docs/VECTOR.md), [InfluxDB](docs/INFLUXDB.md), [alerts](docs/ALERTS.md)

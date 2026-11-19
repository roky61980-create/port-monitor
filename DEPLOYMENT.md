# Ubuntu deployment

Install Go 1.24+, build with `go build -o server ./cmd/server`, install under `/opt/port-monitor`, and configure SQLite writable storage. Install/configure OpenTelemetry Collector using `deploy/otel/otel-collector.yaml`. The systemd unit is in `deploy/systemd/port-monitor.service`; create its service account and adjust database path/permissions before enabling.

Allow TCP/8080 for UI/API, UDP/162 from managed device addresses, and OTLP only on a trusted local/management network. Vector and InfluxDB are optional at this stage: deployment configuration is a starting point and needs a tested record mapping and credentials before production use. Configure SNMPv3 credentials only after v3 trap support has been completed.

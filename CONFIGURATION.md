# Configuration

`config.yaml` supplies bind addresses, poll controls, OTLP destination and SQLite path. Secrets are environment-only: `SNMP_COMMUNITY`, `SNMPV3_AUTH_PASSWORD`, `SNMPV3_PRIV_PASSWORD`, `INFLUXDB_TOKEN`. The current trap listener uses SNMPv2c and defaults to community `public`; replace it through `SNMP_COMMUNITY` in trusted deployments. Do not expose UDP/162 to untrusted networks.

Optical interval defaults to 60 seconds; actual polling is intentionally disabled until a verified provider is implemented. The OTLP endpoint defaults to `127.0.0.1:4317`.

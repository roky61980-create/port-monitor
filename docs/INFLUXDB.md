# Vector → InfluxDB setup note

The current Vector file is a configuration scaffold only; validate it against the deployed Vector version and implement a tested OTLP/JSON-to-line-protocol transform before enabling writes. Supply `INFLUXDB_URL`, `INFLUXDB_ORG`, `INFLUXDB_BUCKET`, and `INFLUXDB_TOKEN` through a protected environment file.

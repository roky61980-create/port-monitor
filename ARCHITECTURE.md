# Architecture

```text
Network device -> UDP SNMP trap receiver -> normalization -> SQLite event/alert state -> alert engine -> SSE dashboard
                                                            -> OpenTelemetry SDK -> OTLP -> Collector
Optical provider interface -> verified sensor readings only -> alerts/telemetry (provider implementation required per device)
```

SQLite holds devices, discovered ports, events, optical latest state and alerts. InfluxDB is reserved for time-series data; this version does not write to InfluxDB. Traps are preferred for port changes; no continuous port polling is implemented. No RMON, CPU, memory or traffic metrics exist.

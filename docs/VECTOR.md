# Vector

`deploy/vector/vector.yaml` is an illustrative HTTP source and InfluxDB sink scaffold. It is not yet a production-ready OTLP-to-Influx conversion: map Collector-exported records to Influx line protocol, test authentication and batching, and validate Vector component compatibility before enabling it.

# OpenTelemetry

Application events are emitted as structured spans with event type, device, vendor, model, interface, alert severity and alert state attributes. The SDK batches OTLP traces to the configured gRPC endpoint. The Collector sample accepts OTLP and uses a debug exporter; configure a tested production exporter/pipeline for downstream systems.

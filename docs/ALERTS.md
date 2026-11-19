# Alerts

Alert states are OPEN, ACKNOWLEDGED and RESOLVED. Repeated down notifications correlate to the active alert by device, interface and event class. Up notifications resolve matching alerts. API actions are `POST /api/alerts/{id}/ack` and `/resolve`; optional acknowledger identity is `X-Actor`. Optical thresholds and hysteresis are not active until verified sensor units and thresholds are provided.

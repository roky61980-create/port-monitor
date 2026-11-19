# Optical DDM

Providers expose `Discover` and `Read` interfaces for Cisco, MikroTik, BDCOM, generic and unknown devices. No vendor OID was verified from repository evidence or supplied MIB/device output in this implementation. Therefore all optical sensor OIDs remain `NOT_VERIFIED`; no readings are returned or fabricated. A provider must not mark a sensor verified until its exact vendor/model/firmware OID, encoding and scale have been verified against vendor documentation/MIB or actual output.

BDCOM is provider-ready only. No RX/TX power, temperature, voltage, bias, module identity or wavelength is claimed. Optical polling and threshold evaluation await a verified provider and configured device credentials.

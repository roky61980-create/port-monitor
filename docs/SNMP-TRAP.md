# SNMP traps

The UDP listener currently accepts SNMPv2c packets using the configured community. It classifies standard IF-MIB linkDown/linkUp notifications and records packet varbinds/source. Device lookup uses sender IP against configured devices; interface text is resolved only when the corresponding port already exists in SQLite. Unknown traps are retained as generic device events.

SNMPv3 trap authentication/decryption, trap queue backpressure policies, device-specific notification mappings, and packet-level integration testing remain future work. Never configure trap sources from untrusted networks.

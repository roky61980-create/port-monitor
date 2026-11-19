# Port alerts

Port/link state changes use traps as the realtime mechanism. Duplicate notifications for an already-open device/interface/type alert update the existing alert instead of creating duplicates. A matching UP event resolves its active alert. Interface discovery and state reconciliation are not yet implemented; populate ports through a verified discovery workflow before depending on interface names.

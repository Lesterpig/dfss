CHANGELOG
=========

v0.1.0
------
> 18/04/2016

#### CLI Client

- Add import command (!35)
- Add export command (!35)
- Add sign command, implement main protocol (!45)
- Add fetch command (!52)
- Ignore hostname verification in TLS (!34)
- Add support for direct IP connection, without using hostnames (!40)
- Add d flag for demonstrator activation (!55)

#### GUI Client

- Add gui binary for client (work in progress, !44)

#### Platform

- Add sequence generation (!36)
- Add signature ignition API (!41)
- Add ttp command (!43)
- Add fetch API (!51)
- Update d flag to use a string instead of boolean (!55)

#### TTP

- Add dfssp binary for TTP (!31, !39)
- Add naive Alert API (!56)

#### Demonstrator

- Add automatic sort of incoming messages (!42)
- Add gui interface for communication visualization (!46)
- Add open/save file (!47)
- Add play/pause/replay actions (!47)
- Add speed selection (!49)
- Add nogui flag (!46)
- Add port flag (!55)

#### Misc

- Use bytes instead of strings for hashes during network communications (!33)
- Add MIT license (!37)
- Use default connection timeout everywhere (30 seconds) (!54)
- Improve deployed binaries size by removing debug symbols

---

v0.0.2
------
> 29/02/2016

- Fix cross-compilation for unix targets
- Fix build error due to grpc upstream update

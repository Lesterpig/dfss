CHANGELOG
=========

v0.3.1
------
> 27/05/2016

#### CLI Client

- Fix unregister command (!94)

#### Platform

- Fix the unability to sign several times the same contract (!93)

v0.3.0
------
> 26/05/2016

#### CLI Client

- Log more information during signature (!69)
- Add an optional local file check before signature (!74)
- Add a slowdown flag for tests (!78)
- Add a stopbefore flag for tests (!79)
- Add TTP resolution (!80)
- Add timeout flag (!81)
- Add unregister command (!85)

#### GUI Client

- Improve contract display (!70)
- Fix crash when the dfss file is corrupted (!71)
- Add an optional local file check before signature (!74)
- Add a cancel button before signature (!75)
- Add TTP resolution (!80)
- Add timeout configuration option (!81)

#### Platform

- Support different subnets between signers (!73)
- Add TTP listing support (!80)
- Disable case sensitivity in mails (!86)
- Add unregister API (!85)
- Fix timeout issues during ready sign (!89)
- Add contract UUID in invitation mails (!90)

#### TTP

- Add TTP resolution (!80)
- Fix concurrency problems (!82)

#### Demonstrator

- Add a step-by-step speed (!76)
- Add communication visualization for TTP and platform (!83)

#### Misc

- Fix the verbose flag, now displays messages sent to demonstrator (!88)

v0.2.0
------
> 25/04/2016

#### CLI Client

- Update command-line interface to match POSIX standards (!64)

#### GUI Client

- Add menu bar (!65)
- Add show contract screen (!65)
- Add fetch contract screen (!65)
- Add signature screen (!61)
- Add basic help message box (!66)
- Add about message box for DFSS (!66)
- Add about message box for Qt (!66)
- Add user mail and selected platform information (!66)
- Improve buttons and feedbacks (!66)
  + Buttons now use the full width of the window
  + Error messages are now displayed in message boxes

#### Platform

- Fix crash if a client disconnects before the ready signal (!60)
- Update command-line interface to match POSIX standards (!63)

#### TTP

- Update command-line interface to match POSIX standards (!64)

#### Demonstrator

- Increase arrow length to 30px (!67)
- Update default unit from nano-second to micro-second (!67)
- Update command-line interface to match POSIX standards (!64)

#### Misc

- Improve security of network communication by checking server's certificate (!62)


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

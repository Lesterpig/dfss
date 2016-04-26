DFSS
====

> Distributed Fair Signing System

DFSS is a simple and secure way to e-sign contracts with a large number of participants, ensuring fairness and minimizing the involvement of the trusted third party (TTP).
This repository contains source code for this INSA Rennes project.

- Website: https://static.lesterpig.com/dfss_web/
- Downloads: https://static.lesterpig.com/dfss/

Project Status
--------------

The DFSS project is developed by fourth year students of the Computer Science department of INSA Rennes (FR).
For now, the infrastructure is working, without the implementation of the signature cryptographic kernel (**Private Contract Signatures**), defined in many research papers, but not yet implemented.

It's thus a *proof of concept* and not production-ready.

Reference:
- [Aybek Mukhamedo, Mark D. Ryan. Fair Multi-party Contract Signing using Private Contract Signatures.](https://www.researchgate.net/publication/222527059_Fair_multi-party_contract_signing_using_private_contract_signatures)
- [Barbara Kordy, Saša Radomirović. Constructing Optimistic Multi-party Contract Signing Protocols.](http://people.irisa.fr/Barbara.Kordy/papers/CSF12.pdf)

Configure workspace
-------------------

1. Install Go (>=1.5) and configure a Go workspace as [explained here](https://golang.org/doc/code.html#Organization)

2. Navigate under `$GOPATH/src` and clone this repository

3. Install build dependencies in `dfss/` directory

```bash
dfss/build/deps.sh
```

4. At this point, you will be able to install the DFSS project with some simple commands

- To install CLI applications:

```bash
go install dfss/dfssc # Client
go install dfss/dfssp # Platform
go install dfss/dfsst # TTP

# or

make install
```

- To build GUI for client into `bin/` directory (using docker image)

```
# You may have to run these commands as root due to docker (sudo won't work)

# Prepare docker image, one time only
make prepare_gui

# Build
make gui
make dfssd
```

Do not attempt to run `go install dfss/...` or `go install ./...`, it won't work due to graphic binaries.

Run dfss modules
----------------

```bash
dfssc help # Client
dfssp help # Platform
dfsst help # TTP
```

For graphic clients, you may need to install some Qt4 libraries on your system.

```bash
cd bin
./gui   # Client GUI
./dfssd # Demonstrator
```

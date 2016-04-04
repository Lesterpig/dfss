DFSS
====

> Distributed Fair Signing System

This repository contains source code for this INSA Rennes project (work in progress).

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
```

Run dfss modules
----------------

```bash
dfssc help # Client
dfssp help # Platform
dfsst help # TTP
```

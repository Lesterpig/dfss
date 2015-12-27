DFSS
====

> Distributed Fair Signing System

This repository contains source code for this INSA Rennes project (work in progress).

Configure workspace
-------------------

1. Install Go (>=1.5) and configure a Go workspace as [explained here](https://golang.org/doc/code.html#Organization)

2. Navigate under `$GOPATH/src` and clone this repository

3. At this point, you will be able to install the DFSS project with a simple command, anywhere from your computer:

```bash
go install dfss/...
```

Run dfss modules
----------------

```bash
dfssc help # Client
dfssp help # Plaform
dfssd help # Demonstrator
```

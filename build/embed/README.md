DFSS
====
> Distributed Fair Signing System

Prerequisites
-------------

DFSS is distributed on the following platforms:

- Linux (amd64, i386 and arm)
- MacOS (amd64)
- Windows (i386)

A running mongoDB database is required for the Platform and the TTP modules.

Installation (UNIX)
-------------------

This archive contains all the needed DFSS modules to run a distributed multiparty signature.
You just have to untar the archive and run the following binaries:

```bash
./dfssc help # Client
./dfssp help # Platform
./dfsst help # TTP
```

On Linux-64 and Windows builds, two additional graphic binaries are included. You may need some Qt4 libraries on your system to run them.

```bash
./dfssc_gui # Graphic client
./dfssd     # Demonstrator
```

### Setup platform and TTP (Trusted Third Party)

The first thing to do is to create the *root certificate of authentication* for the platform.
You can configure several parameters for that (check the `help` command of `dfssp`).

For instance, if we are running the plaform on the `example.com` host:

```bash
./dfssp --cn example.com --country FR --validity 3650 init
```

Then, it is possible to create TTP credentials from generated root credentials.
The generated files are stored in a subdirectory "ttp".

```bash
./dfssp --cn ttp.example.com --country FR --validity 365 ttp
```

You can then start the platform. Here we are considering a mongoDB database running on the same host.
Firstly, we have to configure several environment variables to set smtp server configuration (mails):

```bash
export DFSS_MAIL_SENDER="mailer@example.com"
export DFSS_MAIL_HOST="smtp.example.com"
export DFSS_MAIL_PORT="587"
export DFSS_MAIL_USERNAME="mailer"
export DFSS_MAIL_PASSWORD="password"
```

Then:

```bash
./dfssp start
```

You can also start the TTP:

```bash
./dfsst --cert ttp/cert.pem --key ttp/cert.pem start
```

### Setup clients

Each client needs the `dfssp_rootCA.pem` file in order to connect to the platform in a secure way.
Clients can then register on the platform with the following command:

```bash
./dfssc --ca path/to/dfssp_rootCA.pem --host example.com register
```

A mail will be sent to the user containing a unique token. Use this token to authenticate onto the platform:

```bash
./dfssc --ca path/to/dfssp_rootCA.pem --host example.com auth
```

When this is done, the client will have a certificate and a private key in the current directory.
It's then possible to send new contracts to the platform:

```bash
./dfssc --ca path/to/dfssp_rootCA.pem --host example.com new
```

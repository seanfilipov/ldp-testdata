LDP
===

Copyright (C) 2017-2019 The Open Library Foundation.  This software is 
distributed under the
terms of the Apache License, Version 2.0.  See the file
[LICENSE](https://github.com/folio-org/ldp/blob/master/LICENSE) for
more information.


Overview
--------

The LDP is a database platform to support analytics for
[FOLIO](https://www.folio.org).

**This software is under active development, and no database schema
migrations are currently provided.  Use this software only for testing
purposes.**


System requirements
-------------------

* Linux or macOS
* PostgreSQL 9.6 or later
* [Go](https://golang.org) 1.10 or later


Installing the software
-----------------------

First ensure that the `GOPATH` environment variable specifies a path
that can serve as your Go workspace directory, the place where this
software and other Go packages will be installed.  For example, to set
it to `$HOME/go`:

```shell
$ export GOPATH=$HOME/go
```

Then to download and compile the software:

```shell
$ go get -u github.com/folio-org/ldp/...
```

The compiled executable files, `ldp-update` etc., should appear in `$GOPATH/bin/`.  


Running the LDP
---------------

### Configuration file

Create a configuration file for the LDP:

```ini
# Sample LDP configuration file

[ldp-database]
dbtype = postgres
host = localhost
port = 5432
user = ldpadmin
password = password_goes_here
dbname = ldp
```

The server looks for a configuration file like this one in a location
specified by the `LDP_CONFIG_FILE` environment variable, which
in bash can be set with, for example:

```shell
$ export LDP_CONFIG_FILE=/etc/ldp/ldp.conf
```

### Creating the LDP database

```shell
$ createuser ldpadmin
$ createdb -O ldpadmin ldp
$ ldp-init
```

### Loading data into the database

To load sample data from JSON files in `~/testdata/20181214_043055`:

```shell
$ ldp-update -source ~/testdata/20181214_043055
```

### Schema

The LDP database uses "star schema" with tables prefixed `f_` for fact
and `d_` for dimension, for example, `f_loans` and `d_users`.  It also
includes tables needed for denormalizing new data during updates, for
example, the `groups` table.


### Update process

The `ldp-update` tool performs incremental updates, but it is also used
for batch loads.  It reads one unit of FOLIO data (e.g. a loan
transaction) at a time and adds it to the LDP database, overwriting any
existing data that has the same ID.

If the new data contain foreign key IDs that reference data the LDP also
stores, the update process ensures that the referenced data exist, or if
not, creates placeholder data that can be updated later.  This allows it
to process new data that are streamed from FOLIO "out of order".




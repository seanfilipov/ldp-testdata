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
$ go get -u github.com/folio-org/ldp/cmd/ldp
```

The compiled executable file, `ldp`, should appear in `$GOPATH/bin/`.  


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
dbname = ldpdemo
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
$ createdb -O ldpadmin ldpdemo
$ ldp -init
```

### Loading data into the database

To load sample data from JSON files in `~/testdata/20181214_043055`:

```shell
$ ldp -load -dir ~/testdata/20181214_043055
```

### Schema

The LDP database uses "star schema" and all tables are located in the
`public` schema.  It also includes tables needed for denormalizing new
data during loading, which are stored in the `normal` schema.  The
`loading` schema is for internal use by the data loader.


### Update process

Running the data loader, by using the `-load` flag, performs incremental
loading, but it is also used for batch loads.  It reads one unit of
FOLIO data (e.g. a loan transaction, or a user record) at a time and
adds it to the LDP database.  In the case of "fact" tables, loading
overwrites existing data that have the same FOLIO ID.  In the case of
"dimension" tables, such a collision results in both versions being
preserved and versioned using the `record_effective` attribute.  LDP
primary keys have a `_key` suffix, while FOLIO IDs have a `_id` suffix,
and the primary keys are used to distinguish between different records
with the same FOLIO ID.

If the new data being loaded contain foreign key IDs that reference data
the LDP also stores, the loading process ensures that the referenced
data exist, or if not, creates placeholder data that can be
automatically filled in by a subsequent load containing the needed data.
This will allow the loader to process new data that are streamed from
FOLIO "out of order".



LDP testdata
===

Copyright (C) 2017-2019 The Open Library Foundation.  This software is 
distributed under the
terms of the Apache License, Version 2.0.  See the file
[LICENSE](https://github.com/folio-org/ldp/blob/master/LICENSE) for
more information.


Prerequisites
-------------------

* [Go](https://golang.org) 1.10 or later

Overview
--------

This purpose of this repo is to generate large amounts of fake FOLIO data to support the LDP analytics team.

To download and compile:

```shell
go get github.com/folio-org/ldp-testdata
cd ldp-testdata/cmd/ldp-testdata
go build
cd ../..
```

To run:
```shell
cmd/ldp-testdata/ldp-testdata all
```

Usage
--------
```
cmd/ldp-testdata/ldp-testdata
Usage:
./ldp-testdata FLAGS [all|groups|users|locations|items|loans|circloans|storageitems]
  where FLAGS include:
  -dataFormat string
    	The outputted data format [folioJSON|jsonArray] (default "folioJSON")
  -dir string
    	The directory to use for extract output. If the selected test data depends on
    	other test data (e.g. 'users' depends on 'groups'), that dependency should exist
    	in this directory.
  -nGroups int
    	The number of groups to create (default 12)
  -nItems int
    	The number of items to create (default 10000)
  -nLoans int
    	The number of loans to create (default 10000)
  -nLocations int
    	The number of locations to create (default 20)
  -nUsers int
    	The number of users to create (default 30000)
```

Typically, you will want to run `ldp-testdata all` to generate all data. You can tweak the parameters
using the options, e.g.

```shell
ldp-testdata -nUsers=50000 -nGroups=20 -nLoans=800000 -dir=./myOutput all
```

You can specify the same directory as a previous run to overwrite one type of data:
```shell
ldp-testdata -dir=./myOutput nUsers=20000 users
```

**This software is under active development. Use this software only for testing purposes.**

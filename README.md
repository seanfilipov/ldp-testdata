LDP testdata
===

Copyright (C) 2017-2019 The Open Library Foundation.  This software is 
distributed under the
terms of the Apache License, Version 2.0.  See the file
[LICENSE](https://github.com/folio-org/ldp/blob/master/LICENSE) for
more information.


Overview
--------

This purpose of this repo is to generate large amounts of fake FOLIO data to support the LDP analytics team.

**This software is under active development, and no database schema
migrations are currently provided.  Use this software only for testing
purposes.**


System requirements
-------------------

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
$ go get -u github.com/folio-org/ldp-testdata/cmd/ldp-testdata
```

The compiled executable file, `ldp-testdata`, should appear in `$GOPATH/bin/`.  



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

To download:

```shell
go get github.com/folio-org/ldp-testdata
```

To run:
```shell
go run ./cmd/ldp-testdata/main.go
```

Usage
--------
```
go run ./cmd/ldp-testdata/main.go [FLAGS]

All flags are optional

  -dataFormat string
    	The outputted data format [folioJSON|jsonArray] (default "folioJSON")
  -dir string
    	The directory to store output
  -fileDefs string
    	The filepath of the JSON file definitions (default "filedefs.json")
  -json string
    	JSON array to override the number of objects set filedefs.json
    	Example: '[{"path": "/loan-storage/loans", "n":50000}]'
  -only-json
    	Use with the -json flag to ignore filedefs.json
```

Edit filedefs.json to change the number of objects created for each path, or 
use the `-json` flag to override the number of objects set in filedefs.json

```shell
go run ./cmd/ldp-testdata/main.go -json='[{"path": "/loan-storage/loans", "n":50000}]'
```

**This software is under active development. Use this software only for testing purposes.**

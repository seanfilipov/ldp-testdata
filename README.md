LDP testdata
===

Copyright (C) 2017-2019 The Open Library Foundation.  This software is 
distributed under the
terms of the Apache License, Version 2.0.  See the file
[LICENSE](https://github.com/folio-org/ldp/blob/master/LICENSE) for
more information.


Prerequisites
-------------------

* [Go](https://golang.org) 1.12 or later

Overview
--------

This purpose of this repo is to generate large amounts of fake FOLIO data to support the LDP analytics team.

To download and install:

```shell
go get -u github.com/folio-org/ldp-testdata/...
```

Run the command from the project root:
```shell
cd ~/go/src/folio-org/ldp-testdata
~/go/bin/ldp-testdata
```

Usage
--------
```
~/go/bin/ldp-testdata [FLAGS]

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
  -logLevel string
    	The log level (Trace, Debug, Info, Warning, Error, Fatal and Panic) (default "Info")
  -only-json
    	Use with the -json flag to ignore filedefs.json
```

Edit [filedefs.json](https://github.com/folio-org/ldp-testdata/blob/master/doc/filedefs.md) to change the number of objects created for each path, or 
use the `-json` flag to override the number of objects set in filedefs.json

```shell
~/go/bin/ldp-testdata -json='[{"path": "/loan-storage/loans", "n":50000}]'
```

Supported Routes
--------

- [/groups](https://s3.amazonaws.com/foliodocs/api/mod-users/groups.html)
- [/users](https://s3.amazonaws.com/foliodocs/api/mod-users/users.html)
- [/locations](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/location.html)
- [/location-units/institutions](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/locationunit.html)
- [/service-points](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/service-point.html)
- [/material-types](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/material-type.html)
- [/instance-types](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/instance-type.html)
- [/instance-storage/instances](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/instance-storage.html)
- [/holdings-storage/holdings](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/holdings-storage.html)
- [/item-storage/items](https://s3.amazonaws.com/foliodocs/api/mod-inventory-storage/item-storage.html)
- [/inventory/items](https://s3.amazonaws.com/foliodocs/api/mod-inventory/inventory.html)
- [/loan-storage/loans](https://s3.amazonaws.com/foliodocs/api/mod-circulation-storage/loan-storage.html)
- [/circulation/loans](https://s3.amazonaws.com/foliodocs/api/mod-circulation/circulation.html)

Additional Documentation
--------

- [Algorithms used to generate test data](https://github.com/folio-org/ldp-testdata/blob/master/doc/algorithms.md)
- [Field descriptions for filedefs.json](https://github.com/folio-org/ldp-testdata/blob/master/doc/filedefs.md)

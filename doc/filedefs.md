# filedefs.json

filedefs.json contains the parameters for generating ldp-testdata. An example is below:

```
[
  {
      "module": "mod-users",
      "path": "/groups",
      "objectKey": "usergroups",
      "doc": "https://s3.amazonaws.com/foliodocs/api/mod-users/groups.html",
      "n": 12
  }
]
```

Each object in the array is a `filedef`, defining a file that will be outputed. 

### module

Name for the module that is being simulated

### path

The API path for the data; `path` is the only unique field, so it is treated as the ID for the filedef.

### objectKey

The key for the output JSON file

Example:

```
{
  usergroups: [
    ...
  ]
}
```

### doc

A URL to documention for this API path

### n

The number of objects to generate

For the following filedefs, `n` is irregular:

| Path                | n behavior                             |
|---------------------|----------------------------------------|
| /loan-storage/loans | n is approximate¹                      |
| /inventory/items    | Ignored. Same n as /item-storage/items |
| /circulation/loans  | Ignored. Same n as /loan-storage/loans |

¹ The current simulation for loan objects over 1 year involves some randomness, which makes meeting an exact number difficult.
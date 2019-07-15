# Data Generation Algorithms

In the test data generated so far, I've found that a data generation algorithm can be classified as Simple, Dependent, or Unique.

## Simple

Examples: GenerateGroups(), GenerateLocations(), GenerateMaterialTypes()

Inputs:

- random IDs, e.g. `uuid.NewV4()`
- hardcoded values, e.g. ["dvd", "video recording", "microform", ...]

Output Example:

```
{
    "name": "dvd",
    "id": "39be9a1c-13f7-4f28-8ad1-656282d73a6f"
}
```

## Dependent

Examples: GenerateUsers(), GenerateInstances(), GenerateHoldings()

Inputs: 

- The same inputs as the Simple algorithm
- Data generated from another algorithm, which must be read by file, either by random access or sequentially. In the case of Users, a random patron group ID is selected for each User. In the case of Holdings, all Instances are iterated over so that each Holding has a single parent Instance.

Output Example:

```
{
    "id": "97c0a398-7850-42c0-b0b8-ba600340a691",
    "barcode": "6588112574787113",
    "patronGroup": "b774a12d-3e7e-4fd6-a9a5-ea9170f2548a",
    ...
},
```


## Unique

Example: GenerateLoans()

GenerateLoans() is more complex because of these requirements:

1) The `loanDate` should increment so that the data is over a period of 1 year
2) The `Status` of the loan may be checked in or checked out
3) A maximum of 1 million loans must be supported, which requires that the data be split over multiple files

The algorithm:

- Loop over N, the number of transactions
- Each iteration:
  - Create a loan transaction
  - Increment two counters:
    - the number of transactions to be written to the current file
    - the number of transactions to be written with the current date
  - if N == maxNInFile, write the current loans to file, reset the counter
  - if N == maxNInDay, increment the date, reset the counter

When creating a loan transaction, a random item ID is picked. If that item ID has already been checked out, then the transaction will be to check in that item. Otherwise the item will be checked out. This means that as time progresses, the probability of checkins increases, which is realistic enough. A more realistic algorithm is not a requirement at this time.

Input:
- The same inputs as the Dependent algorithm
- Data that depends on other data in the loop, as described above (date, file, status)

Output Example:

```
{
    "id": "df11cbbd-a664-464f-9ac0-b34394a87222",
    "userId": "3dc73e3b-7e3f-41c5-9342-295f8a308f11",
    "itemId": "e53ff2ff-9875-43a6-af1c-314b484820f6",
    "action": "checkedout",
    "status": {
        "name": "Open"
    },
    "loanDate": "2018-02-06T00:00:00Z",
    "dueDate": "2018-02-20T00:00:00Z"
}
```

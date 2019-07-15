# GenerateLoans()

GenerateLoans() is more complex because of these requirements:
a) the loanDate should increment so that the data is over a period of 1 year
b) the status of the loan may be checked in or checked out
c) a maximum of 1 million loans must be supported, which requires that the data be split over multiple files

- Loop over N, the number of transactions
- Each iteration:
  - Create a loan transaction
  - Increment two counters:
    - the number of transactions to be written to the current file
    - the number of transactions to be written with the current date
  - if N == maxNInFile, write the current loans to file, reset the counter
  - if N == maxNInDay, increment the date, reset the counter

When creating a loan transaction, a random item ID is picked. If that item ID has already been checked out, then the transaction will be to check in that item. Otherwise the item will be checked out. This means that as time progresses, the probability of checkins increases, which is realistic enough. A more realistic algorithm is not a requirement at this time.

# safesql

SafeSQL is a package to provide separated interfaces for master & slave database.

The interfaces are being separated, so the master DB only do the write operations
while the slave only do the read operations.

This package is implemented with very minimal code, the only code needed are:
- open DB
- creating prepare statement

While the other operations like select and query are already inherited from the backend DB which currently use `sqlx`.
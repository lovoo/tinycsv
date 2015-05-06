# tinycsv

A tiny CSV helper toolkit.

## filter

Filter allows you to filter CSV encoded data. It reads either from a file or from stdin.

```
$ filter
filter prints out one or more columns from a CSV and reads from a file or stdin.

Usage of filter:
  -cols="": the column index(es) to be written out to stdout
  -filename="": CSV file (if empty, filter reads from stdin)
  -native=false: If only one column is provided, filter does not escape these line; instead it plainly prints it out.
```

Get it with `go get -u github.com/lovoo/tinycsv/filter` and install it using `go install github.com/lovoo/tinycsv/filter`.

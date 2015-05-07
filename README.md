# tinycsv

tinycsv is a collection of tiny CSV helper tools which are in-use at LOVOO.

Download and installation:

```
go get -u github.com/lovoo/tinycsv/...
go install github.com/lovoo/tinycsv/...
```

## extract

`extract` allows you to extract columns from huge CSV encoded data. It reads either from a file or from stdin and outputs to stdout.

```
$ extract
extract prints out one or more columns from a CSV and reads from a file or stdin.

Usage of ./extract:
  -cols="": the column index(es) to be written out to stdout
  -delim=",": the CSV delimiter; default is ','
  -filename="": CSV file (if empty, which is the default, extract reads from stdin)
  -plain=false: If only one column is provided, extract does not escape these line; instead it plainly prints it out (default false).
  -skipHeader=true: skips the first header line (default true)
```

### Example

```
$ head -n 5 output.csv | extract -cols 3
2015-04-09 19:47:12
2015-04-09 19:36:21
2015-04-09 19:40:22
2015-04-09 19:53:28
2015-04-09 19:56:31
```

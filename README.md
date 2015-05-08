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
$ extract --help
extract prints out one or more columns from a CSV and reads from a file or stdin.

Usage of ./extract:
  -cols="": the column index(es) to be written out to stdout
  -delim=",": the CSV delimiter
  -filename="": CSV file (if empty, which is the default, extract reads from stdin)
  -insertHeader="": inserts a new header line to the output (comma-seperated strings)
  -n=0: Stop after reading n lines (default 0 = unlimited).
  -plain=false: If only one column is provided, extract does not escape these line; instead it plainly prints it out.
  -skipHeader=false: skips the first header line
  -suppress=false: suppress warnings in the input data
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

## summary

`summary` generates and plots a summary for numerical data encoded as CSV. It reads either from a file or from stdin and outputs to stdout.

```
$ summary --help
Usage of ./summary:
  -plot=false: plot a histogram and the standard normal distribution
  -suppress=false: suppress warnings in the input data
```

### Example

```
$ extract -insertHeader foo -cols 0 -filename input.csv -n 1000000 -skipHeader | summary -plot -suppress

      min         max          mean       stddev
foo   -11.510000  647.890000   7.798770   7.843919
```

A generated plot looks like:

![](https://raw.githubusercontent.com/Lovoo/tinycsv/master/histogram.png)

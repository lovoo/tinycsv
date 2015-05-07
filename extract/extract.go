package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func usage() {
	fmt.Printf("extract prints out one or more columns from a CSV and reads from a file or stdin.\n\n")
	flag.Usage()
	os.Exit(1)
}

func main() {
	filename := flag.String("filename", "", "CSV file (if empty, which is the default, extract reads from stdin)")
	cols := flag.String("cols", "", "the column index(es) to be written out to stdout")
	plain := flag.Bool("plain", false, "If only one column is provided, extract does not escape these line; instead it plainly prints it out.")
	delim := flag.String("delim", ",", "the CSV delimiter")
	skipHeader := flag.Bool("skipHeader", false, "skips the first header line")
	insertHeader := flag.String("insertHeader", "", "inserts a new header line to the output (comma-seperated strings)")
	n := flag.Int("n", 0, "Stop after reading n lines (default 0 = unlimited).")
	suppressWarnings := flag.Bool("suppress", false, "suppress warnings in the input data")
	flag.Parse()

	var indexes []int
	maxIndex := 0
	scols := strings.Split(*cols, ",")
	for _, col := range scols {
		if col == "" {
			continue
		}

		i, err := strconv.Atoi(col)
		if err != nil {
			log.Fatal(err)
		}
		indexes = append(indexes, i)

		if maxIndex < i {
			maxIndex = i
		}
	}

	if len(indexes) <= 0 {
		usage()
	}

	r, _ := utf8.DecodeRuneInString(*delim)
	if r == utf8.RuneError {
		log.Fatal("delimiter contains an invalid value.")
	}

	var (
		fd  io.ReadCloser
		err error
	)
	if *filename != "" {
		fd, err = os.Open(*filename)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
	} else {
		fd = os.Stdin
	}

	c := csv.NewReader(bufio.NewReader(fd))
	c.Comma = r

	w := csv.NewWriter(bufio.NewWriter(os.Stdout))
	defer w.Flush()

	line := 0

	if *skipHeader {
		line++
		headers, err := c.Read()
		if err != nil && !*suppressWarnings {
			fmt.Fprintf(os.Stderr, "Could not parse header line (line %d, '%s'): %v\n", line, headers, err)
		}
	}

	if *insertHeader != "" {
		w.Write(strings.Split(*insertHeader, ","))
	}

	for {
		line++
		columns, err := c.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			if !*suppressWarnings {
				fmt.Fprintf(os.Stderr, "Could not parse line (line %d, '%s'): %v\n", line, columns, err)
			}
			continue
		}

		if len(columns) <= maxIndex {
			if !*suppressWarnings {
				fmt.Fprintf(os.Stderr, "Line has not enough columns (line %d): '%s'\n", line, columns)
			}
			continue
		}

		var extractedColumns []string
		for _, i := range indexes {
			extractedColumns = append(extractedColumns, columns[i])
		}
		if *plain && len(extractedColumns) == 1 {
			fmt.Fprintf(os.Stdout, "%s\n", extractedColumns[0])
		} else {
			w.Write(extractedColumns)
		}

		if *n > 0 && line >= *n {
			return
		}
	}
}

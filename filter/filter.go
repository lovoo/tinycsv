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
)

func usage() {
	fmt.Printf("filter prints out one or more columns from a CSV and reads from a file or stdin.\n\n")
	flag.Usage()
	os.Exit(1)
}

func main() {
	filename := flag.String("filename", "", "CSV file (if empty, filter reads from stdin)")
	cols := flag.String("cols", "", "the column index(es) to be written out to stdout")
	plain := flag.Bool("plain", false, "If only one column is provided, filter does not escape these line; instead it plainly prints it out.")
	flag.Parse()

	var indexes []int
	maxIndex := 0
	scols := strings.Split(*cols, ",")
	for _, col := range scols {
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
	w := csv.NewWriter(bufio.NewWriter(os.Stdout))

	line := 0
	for {
		line++
		columns, err := c.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Fprintf(os.Stderr, "Could not parse line (line %d): '%s'\n", line, columns)
			continue
		}

		if len(columns) <= maxIndex {
			fmt.Fprintf(os.Stderr, "Line has not enough columns (line %d): '%s'\n", line, columns)
			continue
		}

		var filteredColumns []string
		for _, i := range indexes {
			filteredColumns = append(filteredColumns, columns[i])
		}
		if *plain && len(filteredColumns) == 1 {
			fmt.Fprintf(os.Stdout, "%s\n", filteredColumns[0])
		} else {
			w.Write(filteredColumns)
		}
	}
}

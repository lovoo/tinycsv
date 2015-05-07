package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/stat"
)

const (
	unknown = iota
	numeric
	text
)

type columnStatistic struct {
	typ  int
	name string
	obsv []float64
	min  float64
	max  float64
}

func usage() {
	fmt.Printf("summary aggregates any input data encoded as CSV and returns key parameters for all desired column indexes.\n\n")
	flag.Usage()
	os.Exit(1)
}

func main() {
	suppressWarnings := flag.Bool("suppress", false, "suppress warnings in the input data")
	doPlot := flag.Bool("plot", false, "plot a histogram and the standard normal distribution")
	flag.Parse()

	c := csv.NewReader(bufio.NewReader(os.Stdin))

	headers, err := c.Read()
	if err != nil {
		log.Fatalf("Could not read header: %v", err)
	}

	stats := make([]*columnStatistic, 0, len(headers))
	for _, name := range headers {
		stats = append(stats, &columnStatistic{
			name: name,
		})
	}

	line := 0
	for {
		line++
		columns, err := c.Read()
		if err != nil {
			if err == io.EOF {
				// Print statistics
				w := tabwriter.NewWriter(os.Stdout, 5, 8, 1, '\t', tabwriter.AlignRight)
				fmt.Fprint(w, "\n\tmin\tmax\tmean\tstddev\n")

				for _, s := range stats {
					mean, stddev := stat.MeanStdDev(s.obsv, nil)
					fmt.Fprintf(w, "%s\t%f\t%f\t%f\t%f\n", s.name, s.min, s.max, mean, stddev)
				}

				w.Flush()

				if *doPlot {
					for _, s := range stats {
						mean, stddev := stat.MeanStdDev(s.obsv, nil)

						p, err := plot.New()
						if err != nil {
							panic(err)
						}
						p.X.Min = s.min - 1.5*stddev
						p.X.Max = s.max + 1.5*stddev
						p.Title.Text = fmt.Sprintf("Histogram for %s", s.name)

						h, err := plotter.NewHist(plotter.Values(s.obsv), 16)
						if err != nil {
							panic(err)
						}
						h.Normalize(1)
						p.Add(h)

						norm := plotter.NewFunction(func(x float64) float64 {
							return 1.0 / (stddev * math.Sqrt(2*math.Pi)) * math.Exp(-((x-mean)*(x-mean))/(2*stddev*stddev))
						})
						norm.Samples = int(p.X.Max-p.X.Min) + 100
						norm.Color = color.RGBA{R: 255, A: 255}
						norm.Width = vg.Points(2)
						p.Add(norm)

						// Save the plot to a PNG file.
						if err := p.Save(4*vg.Inch, 4*vg.Inch, fmt.Sprintf("histogram-%s.png", s.name)); err != nil {
							log.Fatal(err)
						}
					}
				}

				return
			}
			if !*suppressWarnings {
				fmt.Fprintf(os.Stderr, "Could not parse line (line %d, '%s'): %v\n", line, columns, err)
			}
			continue
		}

		for i, data := range columns {
			s := stats[i]

			if s.typ == unknown {
				// Determine the column's data type
				_, err := strconv.ParseFloat(data, 64)
				switch {
				case err == nil:
					// Is numeric
					s.typ = numeric
				case err != nil:
					s.typ = text
				}
			}

			var f float64
			switch s.typ {
			case numeric:
				f, err = strconv.ParseFloat(data, 64)
				if err != nil {
					if !*suppressWarnings {
						fmt.Fprintf(os.Stderr, "Could not parse numeric value (line %d, '%s'): %v\n", line, data, err)
					}
					continue
				}
				if math.IsNaN(f) {
					if !*suppressWarnings {
						fmt.Fprintf(os.Stderr, "Could not parse numeric value (line %d, '%s'): is not a number\n", line, data)
					}
					continue
				}
				if math.IsInf(f, 0) {
					if !*suppressWarnings {
						fmt.Fprintf(os.Stderr, "Could not parse numeric value (line %d, '%s'): infinity\n", line, data)
					}
					continue
				}
			case text:
				// we take the text length as the according metric
				f = float64(len(data))
			default:
				panic("unknown data type")
			}

			s.obsv = append(s.obsv, f)
			if f < s.min {
				s.min = f
			}
			if f > s.max {
				s.max = f
			}
		}
	}
}

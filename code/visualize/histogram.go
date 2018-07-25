// +build OMIT

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var (
	inputDataPath = "../data/cereal.csv"
	outputFolder  = "./output"
)

func main() {
	// Open the advertising dataset file.
	f, err := os.Open(inputDataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	df := dataframe.ReadCSV(f)

	for _, colName := range df.Names() {
		// START OMIT
		plotVals := make(plotter.Values, df.Nrow())
		for i, floatVal := range df.Col(colName).Float() {
			plotVals[i] = floatVal
		}

		h, err := plotter.NewHist(plotVals, 25)
		if err != nil {
			log.Fatal(err)
		}

		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.Title.Text = fmt.Sprintf("Histogram of a %s", colName)

		// Add the histogram to the plot.
		p.Add(h)

		// Save the plot to a file.
		if err := p.Save(4*vg.Inch, 4*vg.Inch, outputFolder+colName+"_hist.png"); err != nil {
			log.Fatal(err)
		}
		// END OMIT
	}
}

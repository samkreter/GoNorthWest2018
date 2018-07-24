package main

import (
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/kniren/gota/dataframe"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var (
	inputDataFile = "../data/cereal.csv"
	outputFolder  = "./output"
)

func main() {

	// Open the advertising dataset file.
	f, err := os.Open(inputDataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a dataframe from the CSV file.
	df := dataframe.ReadCSV(f)

	// Extract the target column.
	yVals := df.Col("rating").Float()

	colNames := []string{"calories", "sodium", "fiber", "carbo", "sugars", "potass", "rating"}

	// Create a scatter plot for each of the features in the dataset.
	for _, colName := range colNames {
		// START OMIT
		pts := make(plotter.XYs, df.Nrow())
		for i, floatVal := range df.Col(colName).Float() {
			pts[i].X = floatVal
			pts[i].Y = yVals[i]
		}

		s, err := plotter.NewScatter(pts)
		if err != nil {
			log.Fatal(err)
		}
		s.GlyphStyle.Color = color.RGBA{R: 233, B: 0, A: 255}
		s.GlyphStyle.Radius = vg.Points(3)

		// Create the plot.
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = "Calories"
		p.Y.Label.Text = "Rating"
		p.Add(plotter.NewGrid())

		// Save the plot to a PNG file.
		p.Add(s)
		err = p.Save(4*vg.Inch, 4*vg.Inch, filepath.Join(outputFolder, colName+"_scatter.png"))
		if err != nil {
			log.Fatal(err)
		}
		// END OMIT
	}
}

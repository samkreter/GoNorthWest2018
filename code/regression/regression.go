package main

import (
	"encoding/csv"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/kniren/gota/dataframe"
	"github.com/sajari/regression"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var (
	inputDataPath = "../data/training.csv"
	outputFolder  = ""
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	//multi()
	//single()
	plotRegression()
}

func single() {
	// Open the training dataset file.
	f, err := os.Open(inputDataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader reading from the opened file.
	reader := csv.NewReader(f)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 3
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// START OMIT
	var r regression.Regression
	r.SetObserved("rating")
	r.SetVar(0, "Sugars")

	// Loop over the CSV records adding the training data.
	for i, record := range trainingData {
		// Skip the header.
		if i == 0 {
			continue
		}

		ratingVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}
		sugarsVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Add these points to the regression value.
		r.Train(regression.DataPoint(ratingVal, []float64{sugarsVal}))
	}
	// END OMIT

	// Train the regression model.
	r.Run()

	// Output the trained model
	fmt.Printf("\nModel Info:\n%v\n\n", r.Formula)

	// Open the training dataset file.
	testFile, err := os.Open("../data/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer testFile.Close()

	// Create a new CSV reader reading from the opened file.
	reader = csv.NewReader(testFile)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 3
	testData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Loop over the test data predicting y and evaluating the prediction
	// with the mean absolute error.
	var mAE float64
	for i, record := range testData {

		// Skip the header.
		if i == 0 {
			continue
		}

		// Parse the observed diabetes progression measure, or "y".
		yObserved, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the bmi value.
		sugarVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Predict y with our trained model.
		yPredicted, err := r.Predict([]float64{sugarVal})

		// Add the to the mean absolute error.
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	// Output the MAE to standard out.
	fmt.Printf("MAE = %0.2f\n\n", mAE)
}

func multi() {
	// Open the training dataset file.
	f, err := os.Open(inputDataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new CSV reader reading from the opened file.
	reader := csv.NewReader(f)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 3
	trainingData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// In this case we are going to try and model our Sales
	// by the TV and Radio features plus an intercept.
	var r regression.Regression
	r.SetObserved("rating")
	r.SetVar(0, "Sugars")
	r.SetVar(1, "Calories")

	// Loop over the CSV records adding the training data.
	for i, record := range trainingData {

		// Skip the header.
		if i == 0 {
			continue
		}

		ratingVal, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		sugarsVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		caloriesVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Add these points to the regression value.
		r.Train(regression.DataPoint(ratingVal, []float64{sugarsVal, caloriesVal}))
	}

	// Train the regression model.
	r.Run()

	// Output the trained model
	fmt.Printf("\nModel Info:\n%v\n\n", r.Formula)

	// Open the training dataset file.
	testFile, err := os.Open("../data/test.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer testFile.Close()

	// Create a new CSV reader reading from the opened file.
	reader = csv.NewReader(testFile)

	// Read in all of the CSV records
	reader.FieldsPerRecord = 3
	testData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Loop over the test data predicting y and evaluating the prediction
	// with the mean absolute error.
	var mAE float64
	for i, record := range testData {

		// Skip the header.
		if i == 0 {
			continue
		}

		// Parse the observed diabetes progression measure, or "y".
		yObserved, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Parse the bmi value.
		sugarVal, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		caloriesVal, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			log.Fatal(err)
		}

		// Predict y with our trained model.
		yPredicted, err := r.Predict([]float64{sugarVal, caloriesVal})

		// Add the to the mean absolute error.
		mAE += math.Abs(yObserved-yPredicted) / float64(len(testData))
	}

	// Output the MAE to standard out.
	fmt.Printf("MAE = %0.2f\n\n", mAE)
}

func plotRegression() {
	f, err := os.Open(inputDataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a dataframe
	df := dataframe.ReadCSV(f)

	// Extract the target column.
	yVals := df.Col("rating").Float()

	// pts will hold the values for plotting.
	pts := make(plotter.XYs, df.Nrow())

	// ptsPred will hold the predicted values for plotting.
	ptsPred := make(plotter.XYs, df.Nrow())

	// Fill pts with data.
	for i, floatVal := range df.Col("sugars").Float() {
		pts[i].X = floatVal
		pts[i].Y = yVals[i]
		ptsPred[i].X = floatVal
		ptsPred[i].Y = predict(floatVal)
	}

	// Create the plot.
	p, err := plot.New()
	if err != nil {
		log.Fatal(err)
	}
	p.X.Label.Text = "Sugars"
	p.Y.Label.Text = "Rating"
	p.Add(plotter.NewGrid())

	// Add the scatter plot points for the observations.
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}
	s.GlyphStyle.Color = color.RGBA{R: 233, B: 0, A: 255}
	s.GlyphStyle.Radius = vg.Points(3)

	// Add the line plot points for the predictions.
	l, err := plotter.NewLine(ptsPred)
	if err != nil {
		log.Fatal(err)
	}
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	// Save the plot to a PNG file.
	p.Add(s, l)
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "regression_line.png"); err != nil {
		log.Fatal(err)
	}
}

// predict uses our trained regression model to made a prediction.
func predict(sugar float64) float64 {
	//return 89.38 + sugar*-0.45 // Calories
	return 58.83 + sugar*-2.36 // Sugars
}

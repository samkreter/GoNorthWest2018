package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/sajari/regression"
)

// ModelInfo includes the information about the
// model that is output from the training.
type ModelInfo struct {
	Intercept    float64           `json:"intercept"`
	Coefficients []CoefficientInfo `json:"coefficients"`
}

// CoefficientInfo include information about a
// particular model coefficient.
type CoefficientInfo struct {
	Name        string  `json:"name"`
	Coefficient float64 `json:"coefficient"`
}

var inFile string
var outDir string

func main() {

	flag.StringVar(&inFile, "inFile", "../data/training.csv", "The file with the training data.")
	flag.StringVar(&outDir, "outDir", "", "The output directory")

	flag.Parse()

	// Open the training dataset file.
	f, err := os.Open(inFile)
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

		// ignoring the header.
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

	// Fill in the model information.
	modelInfo := ModelInfo{
		Intercept: r.Coeff(0),
		Coefficients: []CoefficientInfo{
			CoefficientInfo{
				Name:        "sugars",
				Coefficient: r.Coeff(1),
			},
			CoefficientInfo{
				Name:        "calories",
				Coefficient: r.Coeff(2),
			},
		},
	}

	// Marshal the model information.
	outputData, err := json.MarshalIndent(modelInfo, "", "    ")
	if err != nil {
		log.Fatal(err)
	}

	// Save the marshalled output to a file.
	if err := ioutil.WriteFile(filepath.Join(outDir, "model.json"), outputData, 0644); err != nil {
		log.Fatal(err)
	}
}

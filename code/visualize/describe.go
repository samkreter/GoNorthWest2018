package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kniren/gota/dataframe"
)

var (
	inputDataPath = "../data/cereal.csv"
	outputFolder  = ""
)

func main() {
	f, err := os.Open(inputDataPath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a dataframe
	df := dataframe.ReadCSV(f)

	fmt.Println(df.Describe())
}

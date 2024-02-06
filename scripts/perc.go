package main

import (
	"encoding/json"
	"fmt"
	"glidingBox/services"
	"glidingBox/src"
	"io/fs"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

type ResultLine struct {
	Name    string                     `json:"name"`
	Parent  string          		   `json:"folder"`
	Results []src.LocalPercolationData `json:"percolation_results"`
}

func ProcessAnImage(folderPath string, file fs.FileInfo, kernelStart int,  kernelEnd int) ResultLine {
	kernelNumber := (kernelEnd-kernelStart)/2 + 1

	filePath := folderPath + file.Name()
	img, err := services.OpenImage(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var results []src.LocalPercolationData

	sem := make(chan src.LocalPercolationData, kernelNumber)

	image := src.EncodeInterleave(img)
	for i := 0; i < kernelNumber; i++ {
		go func(i, kernelStart int) {
			// Process each kernel in parallel.
			kernelSize := kernelStart + i*2
			probabilities := src.PercolationSimple(image, kernelSize)
			sem <- probabilities
		}(i, kernelStart)
	}

	fmt.Print("Concluídos: ")

	for i := 0; i < kernelNumber; i++ {
		r := <-sem
		fmt.Printf("%d ", r.KernelSize)
		results = append(results, r)
	}
	fmt.Print("\n")
	return ResultLine{
		Name: file.Name(),
		Parent: folderPath,
		Results: results ,
	}
}

func processDir(path string, kernelStart int,  kernelEnd int, processChildren bool) []ResultLine {
	runtime.GOMAXPROCS(4)
	
	// Read all files inside input dir
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	
	totalFiles := len(files)
	fmt.Printf("Processing %s (%d files)\n", path, totalFiles)
	
	// Iterate over all images
	var results []ResultLine
	for idx, f := range files {
		if f.IsDir() {
			if !processChildren {
				continue
			}
			childPath := path + f.Name() + "/"
			childResult := processDir(childPath, kernelStart, kernelEnd, processChildren)
			results = append(results, childResult...)

		} else {
			fmt.Printf("[%d/%d] %s\n", idx+1, totalFiles, f.Name())
			start := time.Now()
			results = append(results, ProcessAnImage(path, f, kernelStart, kernelEnd))
			fmt.Printf("Elapsed time %v\n", time.Now().Sub(start))
		}
	}

	return results
}

func SaveResults(results []ResultLine, filename string) {
	// Export to json
	b, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return
	}

	permissions := 0644
	err = ioutil.WriteFile(filename, b, fs.FileMode(permissions))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	inputDir := "Full-folder-path/"
	outputFile := "percolation-data.json"

	results := processDir(inputDir, 3, 41, true)

	SaveResults(results, outputFile)
}

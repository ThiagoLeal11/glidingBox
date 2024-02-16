package main

import (
	"encoding/json"
	"fmt"
	"glidingBox/services"
	"glidingBox/src"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"runtime"
	"time"
)

type ResultLine struct {
	Name    string                     `json:"name"`
	Parent  string                     `json:"folder"`
	Results []src.LocalPercolationData `json:"percolation_results"`
}

func processAnImage(folderPath string, file fs.FileInfo, kernelStart int, kernelEnd int) ResultLine {
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
		Name:    file.Name(),
		Parent:  folderPath,
		Results: results,
	}
}

func processBatch(batch []fs.FileInfo, folderPath string, kernelStart int, kernelEnd int, includeChildren bool) []ResultLine {
	var results []ResultLine
	batchSize := len(batch)
	for idx, f := range batch {
		fmt.Printf("[%d/%d] %s\n", idx+1, batchSize, f.Name())
		if f.IsDir() {
			if !includeChildren {
				continue
			}
			childPath := fmt.Sprintf("%s%s/", folderPath, f.Name())
			childResult := processChild(childPath, kernelStart, kernelEnd)
			results = append(results, childResult...)
			continue
		}

		start := time.Now()
		results = append(results, processAnImage(folderPath, f, kernelStart, kernelEnd))
		fmt.Printf("Elapsed time %v\n", time.Since(start))
	}
	return results
}

func processChild(path string, kernelStart int, kernelEnd int) []ResultLine {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	totalFiles := len(files)
	fmt.Printf("Processing %s (%d files)\n", path, totalFiles)
	return processBatch(files, path, kernelStart, kernelEnd, true)
}

func processDir(path string, kernelStart int, kernelEnd int, processChildren bool, batchSize int, outputFile string, includeChildren bool) {

	runtime.GOMAXPROCS(4)

	// Read all files inside input dir
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	totalFiles := len(files)
	fmt.Printf("Processing %s (%d files)\n", path, totalFiles)

	batches := [][]fs.FileInfo{files}
	if batchSize > 0 {
		batches = getBatches(files, batchSize)
		fmt.Printf("Split process in %d batches\n", len(batches))
	}
	batchCount := len(batches)

	for idx, batch := range batches {
		outputName := outputFile

		if batchCount > 1 {
			fmt.Printf("Processing batch %d of %d\n", idx+1, batchCount)
			outputName = fmt.Sprintf("(%d of %d) %s", idx+1, batchCount, outputFile)
		}

		result := processBatch(batch, path, kernelStart, kernelEnd, includeChildren)
		saveResults(result, outputName)
	}
}

func saveResults(results []ResultLine, filename string) {
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

func getBatches(arr []fs.FileInfo, batchSize int) [][]fs.FileInfo {
	total := len(arr)
	batchCount := int(math.Ceil(float64(total) / float64(batchSize)))
	var batches [][]fs.FileInfo

	for idx := 0; idx < batchCount; idx++ {
		startingIdx := idx * batchSize
		endingIdx := startingIdx + batchSize
		if endingIdx > total {
			endingIdx = total
		}
		batches = append(batches, arr[startingIdx:endingIdx])
	}

	return batches
}

func main() {

	inputDir := "Full-folder-path/"
	outputFile := "percolation-data.json"

	processDir(inputDir, 3, 41, true, 2, outputFile, true)
}

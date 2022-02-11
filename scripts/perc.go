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
	Results []src.LocalPercolationData `json:"percolation_results"`
}

func ProcessAnImage(inputDir string, file fs.FileInfo) []src.LocalPercolationData {
	fmt.Printf("Processando imagem %s\n", file.Name())
	kernelStart := 3
	kernelEnd := 41
	kernelNumber := (kernelEnd-kernelStart)/2 + 1

	filePath := inputDir + file.Name()
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
	return results
}

func processDir(inputDir string, outputFile string){
	runtime.GOMAXPROCS(4)
	
	// Read all files inside input dir
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	
	// Iterate over all images
	var results []ResultLine
	for _, f := range files {
		start := time.Now()
		results = append(results, ResultLine{
			Name:    f.Name(),
			Results: ProcessAnImage(inputDir, f),
		})
		fmt.Printf("Elapsed time %v\n", time.Now().Sub(start))
	}

	// Export to json
	b, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
		return
	}

	permissions := 0644
	err = ioutil.WriteFile(outputFile, b, fs.FileMode(permissions))
	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	inputDir := "Full-folder-path"
	outputFile := "percolation-data.json"

	processDir(inputDir,outputFile)		
}

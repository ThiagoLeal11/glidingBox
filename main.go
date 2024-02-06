package main

import (
	"encoding/json"
	"fmt"
	"glidingBox/services"
	"glidingBox/src"
	"io/fs"
	"io/ioutil"
	"log"
)

type ResultKernels struct {
	Probabilities []float64 `json:"probabilities"`
	KernelSize    int       `json:"kernel_size"`
}

type ResultLine struct {
	Name    string          `json:"name"`
	Parent  string          `json:"folder"`
	Results []ResultKernels `json:"result_kernels"`
}

func ProcessAnImage(folderPath string, file fs.FileInfo, kernelStart int,  kernelEnd int) ResultLine {
	kernelNumber := (kernelEnd-kernelStart)/2 + 1

	path := folderPath + file.Name()
	img, err := services.OpenImage(path)
	if err != nil {
		log.Fatal(err)
	}

	var results []ResultKernels

	sem := make(chan ResultKernels, kernelNumber)

	image := src.EncodeInterleave(img)
	for i := 0; i < kernelNumber; i++ {
		go func(i, kernelStart int) {
			kernelSize := kernelStart + i*2
			probabilities := src.GlidingBoxSimple(image, kernelSize)
			sem <- ResultKernels{
				Probabilities: probabilities,
				KernelSize:    kernelSize,
			}
		}(i, kernelStart)
	}

	fmt.Print("Concluídos: ")

	for i := 0; i < kernelNumber; i++ {
		r := <-sem
		fmt.Printf("%d ", r.KernelSize)
		results = append(results, r)
	}
	fmt.Print("\n")
	return ResultLine {
		Name: file.Name(),
		Parent: folderPath,
		Results: results ,
	}
}

func processDir(path string, kernelStart int,  kernelEnd int, processChildren bool) []ResultLine {
	
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
			results = append(results, ProcessAnImage(path, f, kernelStart, kernelEnd))
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
	outputFile := "probability-matrix.json"

	results := processDir(inputDir, 3, 41, true)

	SaveResults(results, outputFile)

}

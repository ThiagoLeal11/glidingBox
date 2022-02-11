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
	Results []ResultKernels `json:"result_kernels"`
}

func ProcessAnImage(inputDir string, file fs.FileInfo) []ResultKernels {
	fmt.Printf("Processango imagem %s\n", file.Name())
	kernelStart := 3
	kernelEnd := 41
	kernelNumber := (kernelEnd-kernelStart)/2 + 1

	filePath := inputDir + file.Name()
	img, err := services.OpenImage(filePath)
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
	return results
}

func processDir(inputDir string, outputFileName string){
	
	outputFile := inputDir + outputFileName

	// Read all files inside input dir
	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate over all images
	var results []ResultLine
	for _, f := range files {
		results = append(results, ResultLine{
			Name:    f.Name(),
			Results: ProcessAnImage(inputDir, f),
		})
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
	outputFile := "probability-matrix.json"

	processDir(inputDir,outputFile)

	
}

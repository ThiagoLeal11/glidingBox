package main

import (
	"fmt"
	"glidingBox/services"
	"glidingBox/services/functions"
	"time"
)

func main() {
	img, err := services.OpenImage("/home/thiago/go/src/glidingBox/assets/Benign (1).png")
	if err != nil {
		return
	}

	start := time.Now()
	image := services.EncodeInterleave(img)

	kernelStart := 3
	kernelEnd := 3
	kernelNumber := (kernelEnd-kernelStart)/2 + 1

	//data := make([]int32, kernelNumber)
	//sem := make(chan int, kernelNumber)

	// Make in parallel
	for i := 0; i < kernelNumber; i++ {
		//go func(kernelIdx, kernelStart int) {
		kernelSize := kernelStart + i*2
		fmt.Println(kernelSize)
		results := functions.GlidingBox(image, kernelSize)
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println(elapsed)
		baseResult := functions.GlidingBoxSimple(image, kernelSize)
		//sem <- 1
		fmt.Println(results)
		fmt.Println(baseResult)
		//}(i, kernelStart)
	}

	// Wait for goroutines to finish
	//for i := 0; i < kernelNumber; i++ {
	//	<- sem
	//}

	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
}

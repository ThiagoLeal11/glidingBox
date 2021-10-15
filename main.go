package main

import (
	"fmt"
	"glidingBox/services"
	"time"
)

func main() {
	img, err := services.OpenImage("/home/thiago/go/src/glidingBox/assets/Benign (1).png")
	if err!=nil{
		return
	}

	start := time.Now()
	image := services.EncodeInterleave(img)
	for i := 3; i < 45; i += 2 {
		fmt.Println(i)
		results := services.GlidingBox(image, i)
		fmt.Println(results[0], results[1], results[2], results[3])
	}
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println(elapsed)
	//pixel := img.At(1, 0)
	//
	//model := color.NRGBAModel
	//result := model.Convert(pixel)
	//
	//fmt.Println("sim:")
	//fmt.Println(img.Bounds().Size().X, img.Bounds().Size().Y)
	//fmt.Println(result)
	//
	//
	//fmt.Println(pixel)
	//r, g, b, a := result.RGBA()
	//fmt.Println(uint8(r), uint8(g), uint8(b), uint8(a))


}

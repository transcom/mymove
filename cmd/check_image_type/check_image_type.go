package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	img := flag.String("image", "", "path to image to check the type of")
	flag.Parse()

	if len(*img) == 0 {
		log.Fatal("Please provide an image to check using the -image flag")
	}
	imgBytes, err := ioutil.ReadFile(*img)
	if err != nil {
		log.Fatal(err)
	}
	contentType := http.DetectContentType(imgBytes)
	fmt.Println(contentType)
}

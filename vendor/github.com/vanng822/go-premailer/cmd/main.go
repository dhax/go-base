package main

import (
	"flag"
	"fmt"
	"github.com/vanng822/go-premailer/premailer"
	"log"
	"os"
	"time"
)

func main() {
	var (
		inputFile           string
		outputFile          string
		removeClasses       bool
		skipCssToAttributes bool
	)
	flag.StringVar(&inputFile, "i", "", "Input file")
	flag.StringVar(&outputFile, "o", "", "Output file")
	flag.BoolVar(&removeClasses, "remove-classes", false, "Remove class attribute")
	flag.BoolVar(&skipCssToAttributes, "skip-css-to-attributes", false, "No copy of css property to html attribute")
	flag.Parse()
	if inputFile == "" {
		flag.Usage()
		return
	}
	start := time.Now()
	options := premailer.NewOptions()
	options.RemoveClasses = removeClasses
	options.CssToAttributes = !skipCssToAttributes
	prem := premailer.NewPremailerFromFile(inputFile, options)
	html, err := prem.Transform()
	log.Printf("took: %v", time.Now().Sub(start))
	if err != nil {
		log.Fatal(err)
	}
	if outputFile != "" {
		fd, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer fd.Close()
		fd.WriteString(html)
	} else {
		fmt.Println(html)
	}
}

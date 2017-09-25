package premailer

import (
	"github.com/PuerkitoBio/goquery"
	"os"
)

// NewPremailerFromFile take an filename
// Read the content of this file
// and create a goquery.Document
// and then create and Premailer instance.
// It will panic if any error happens
func NewPremailerFromFile(filename string, options *Options) Premailer {
	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fd.Close()
	d, err := goquery.NewDocumentFromReader(fd)
	if err != nil {
		panic(err)
	}
	return NewPremailer(d, options)
}

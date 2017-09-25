package premailer

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
)

// NewPremailerFromString take in a document in string format
// and create a goquery.Document
// and then create and Premailer instance.
// It will panic if any error happens
func NewPremailerFromString(doc string, options *Options) Premailer {
	read := strings.NewReader(doc)
	d, err := goquery.NewDocumentFromReader(read)
	if err != nil {
		panic(err)
	}
	return NewPremailer(d, options)
}

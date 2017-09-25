package premailer

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSupportedSelectors(t *testing.T) {
	selectors := []string{
		".footer__content_wrapper--last",
		"table[class=\"body\"] .footer__content td",
		"table[class=\"body\"] td.footer__link_wrapper--first",
		".header + .content",
		"#firstname",
		"p ~ ul",
		"div > p",
		"div > p",
		"div p",
		"div, p",
		"[target]",
		"[target=_blank]",
		"[title~=flower]",
		"[lang|=en]",
		"a[href^=\"https\"]",
		"a[href$=\".pdf\"]",
		"a[href*=\"css\"]",
		"p:empty",
		"p:first-child",
		"p:first-of-type",
		"p:last-child",
		"p:last-of-type",
		":not(p)",
		"p:nth-child(2)",
		"p:nth-last-child(2)",
		"p:nth-of-type(2)",
		"p:only-child",
		"p:nth-last-of-type(2)",
		"div:not(:nth-child(1))",
		"div:not(:not(:first-child))",
	}

	html := `<html>
        <head>
        <title>Title</title>
        <style type="text/css" data-premailer="ignore">
        h1, h2 {
        	color:red;
        }
        strong {
        	text-decoration:none
        }
        </style>
        </head>
        <body>
        <h1>Hi!</h1>
        <p><strong>Yes!</strong></p>
        </body>
        </html>`

	read := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(read)

	pr := premailer{}
	pr.doc = doc
	for _, selector := range selectors {
		assert.NotPanics(t, func() {
			pr.doc.Find(selector)
		})
	}
}

func TestNotSupportedSelectors(t *testing.T) {

	html := `<html>
        <head>
        <title>Title</title>
        <style type="text/css" data-premailer="ignore">
        h1, h2 {
        	color:red;
        }
        strong {
        	text-decoration:none
        }
        </style>
        </head>
        <body>
        <h1>Hi!</h1>
        <p><strong>Yes!</strong></p>
        </body>
        </html>`

	read := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(read)

	pr := premailer{}
	pr.doc = doc

	notSupportedSelectors := []string{
		"input:checked",
		"input:disabled",
		"input:enabled",
		"input:optional",
		"input:read-only",
		"p:lang(it)",
	}

	for _, selector := range notSupportedSelectors {
		assert.Equal(t, 0, pr.doc.Find(selector).Length())
	}
}

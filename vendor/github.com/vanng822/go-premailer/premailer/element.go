package premailer

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/vanng822/css"
	"sort"
	"strings"
)

type elementRules struct {
	element         *goquery.Selection
	rules           []*styleRule
	cssToAttributes bool
}

func (er *elementRules) inline() {
	inline, _ := er.element.Attr("style")

	var inlineStyles map[string]*css.CSSStyleDeclaration
	if inline != "" {
		inlineStyles = css.ParseBlock(inline)
	}

	styles := make(map[string]string)
	for _, rule := range er.rules {
		for prop, s := range rule.styles {
			styles[prop] = s.Value
		}
	}

	if len(inlineStyles) > 0 {
		for prop, s := range inlineStyles {
			styles[prop] = s.Value
		}
	}

	final := make([]string, 0, len(styles))
	for p, v := range styles {
		final = append(final, fmt.Sprintf("%s:%s", p, v))
		if er.cssToAttributes {
			er.style_to_basic_html_attribute(p, v)
		}
	}

	sort.Strings(final)
	style := strings.Join(final, ";")
	if style != "" {
		er.element.SetAttr("style", style)
	}

}

func (er *elementRules) style_to_basic_html_attribute(prop, value string) {
	switch prop {
	case "text-align":
		er.element.SetAttr("align", value)
	case "vertical-align":
		er.element.SetAttr("valign", value)
	case "background-color":
		er.element.SetAttr("bgcolor", value)
	case "width":
		fallthrough
	case "height":
		if strings.HasSuffix(value, "px") {
			value = value[:len(value)-2]
		}
		er.element.SetAttr(prop, value)
	}
}

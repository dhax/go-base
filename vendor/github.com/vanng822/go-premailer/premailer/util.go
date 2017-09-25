package premailer

import (
	"github.com/vanng822/css"
)

func copyRule(selector string, rule *css.CSSRule) *css.CSSRule {
	// copy rule for each selector
	styles := make(map[string]*css.CSSStyleDeclaration)
	for prop, s := range rule.Style.Styles {
		styles[prop] = css.NewCSSStyleDeclaration(s.Property, s.Value, s.Important)
	}
	copiedStyle := css.CSSStyleRule{SelectorText: selector, Styles: styles}
	copiedRule := &css.CSSRule{Type: rule.Type, Style: copiedStyle}
	return copiedRule
}

func makeRuleImportant(rule *css.CSSRule) string {
	// this for using Text() which has nice sorted props
	for _, s := range rule.Style.Styles {
		s.Important = 1
	}
	return rule.Style.Text()
}

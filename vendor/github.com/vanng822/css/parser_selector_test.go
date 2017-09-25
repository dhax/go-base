package css

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMultipleSelectors(t *testing.T) {
	css := Parse(`div .a {
						font-size: 150%;
					}
					p .b {
						font-size: 250%;
					}`)

	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "div .a")
	assert.Equal(t, css.CssRuleList[1].Style.SelectorText, "p .b")

}

func TestIdSelector(t *testing.T) {
	css := Parse("#div { color: red;}")

	assert.Equal(t, css.CssRuleList[0].Style.Styles["color"].Value, "red")
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "#div")
}

func TestClassSelector(t *testing.T) {
	css := Parse(".div { color: green;}")

	assert.Equal(t, css.CssRuleList[0].Style.Styles["color"].Value, "green")
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, ".div")
}

func TestStarSelector(t *testing.T) {
	css := Parse("* { text-rendering: optimizelegibility; }")

	assert.Equal(t, "optimizelegibility", css.CssRuleList[0].Style.Styles["text-rendering"].Value)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "*")
}

func TestStarSelectorMulti(t *testing.T) {
	css := Parse(`div .a {
						font-size: 150%;
					}
				* { text-rendering: optimizelegibility; }`)

	assert.Equal(t, "150%", css.CssRuleList[0].Style.Styles["font-size"].Value)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "div .a")

	assert.Equal(t, "optimizelegibility", css.CssRuleList[1].Style.Styles["text-rendering"].Value)
	assert.Equal(t, css.CssRuleList[1].Style.SelectorText, "*")
}

func TestMixedClassSelectors(t *testing.T) {
	selectors := []string{".footer__content_wrapper--last",
		"table[class=\"body\"] .footer__content td",
		"table[class=\"body\"] td.footer__link_wrapper--first",
		"table[class=\"body\"] td.footer__link_wrapper--last"}

	for _, selector := range selectors {
		css := Parse(fmt.Sprintf(` %s {
							    	border-collapse: separate;
							    	padding: 10px 0 0
							    	}`, selector))

		assert.Equal(t, "separate", css.CssRuleList[0].Style.Styles["border-collapse"].Value)
		assert.Equal(t, "10px 0 0", css.CssRuleList[0].Style.Styles["padding"].Value)
		assert.Equal(t, css.CssRuleList[0].Style.SelectorText, selector)
	}
}

func TestGenericSelectors(t *testing.T) {
	selectors := []string{
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
		".header + .content",
		"#firstname",
		"table[class=\"body\"] .footer__content td",
		"table[class=\"body\"] td.footer__link_wrapper--first"}

	for _, selector := range selectors {
		css := Parse(fmt.Sprintf(` %s {
							    	border-collapse: separate;
							    	padding: 10px 0 0
							    	}`, selector))

		assert.Equal(t, "separate", css.CssRuleList[0].Style.Styles["border-collapse"].Value)
		assert.Equal(t, "10px 0 0", css.CssRuleList[0].Style.Styles["padding"].Value)
		assert.Equal(t, css.CssRuleList[0].Style.SelectorText, selector)
	}
}

func TestFilterSelectors(t *testing.T) {
	selectors := []string{
		"a:active",
		"p::after",
		"p::before",
		"input:checked",
		"input:disabled",
		"input:in-range",
		"input:invalid",
		"input:optional",
		"input:read-only",
		"input:enabled",
		"p:empty",
		"p:first-child",
		"p::first-letter",
		"p::first-line",
		"p:first-of-type",
		"input:focus",
		"a:hover",
		"p:lang(it)",
		"p:last-child",
		"p:last-of-type",
		"a:link",
		":not(p)",
		"p:nth-child(2)",
		"p:nth-last-child(2)",
		"p:only-of-type",
		"p:only-child",
		"p:nth-last-of-type(2)",
		"div:not(:nth-child(1))",
		"div:not(:not(:first-child))",
		":root",
		"::selection",
		"#news:target"}

	for _, selector := range selectors {
		css := Parse(fmt.Sprintf(` %s {
							    	border-collapse: separate;
							    	padding: 10px 0 0
							    	}`, selector))

		assert.Equal(t, "separate", css.CssRuleList[0].Style.Styles["border-collapse"].Value)
		assert.Equal(t, "10px 0 0", css.CssRuleList[0].Style.Styles["padding"].Value)
		assert.Equal(t, css.CssRuleList[0].Style.SelectorText, selector)
	}
}

func TestFontFace(t *testing.T) {
	css := Parse(`@font-face {
				      font-family: "Bitstream Vera Serif Bold";
				      src: url("https://mdn.mozillademos.org/files/2468/VeraSeBd.ttf");
				    }
				    
				    body { font-family: "Bitstream Vera Serif Bold", serif }`)

	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-family"].Value, "\"Bitstream Vera Serif Bold\"")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["src"].Value, "url(\"https://mdn.mozillademos.org/files/2468/VeraSeBd.ttf\")")
	assert.Equal(t, css.CssRuleList[1].Style.Styles["font-family"].Value, "\"Bitstream Vera Serif Bold\", serif")
	assert.Equal(t, css.CssRuleList[0].Type, FONT_FACE_RULE)
}

func TestPage(t *testing.T) {
	css := Parse(`@page :first {
					margin: 2in 3in;
				}`)

	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, ":first")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["margin"].Value, "2in 3in")
	assert.Equal(t, css.CssRuleList[0].Type, PAGE_RULE)
}

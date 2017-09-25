package css

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithoutImpotant(t *testing.T) {
	css := Parse(`div .a { font-size: 150%;}`)
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Value, "150%")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Property, "font-size")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Important, 0)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "div .a")

}

func TestWithImpotant(t *testing.T) {
	css := Parse("div .a { font-size: 150% !important;}")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Value, "150%")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Property, "font-size")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Important, 1)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "div .a")
}

func TestMultipleDeclarations(t *testing.T) {
	css := Parse(`div .a {
				font-size: 150%;
				width: 100%
				}`)
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Value, "150%")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Property, "font-size")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Important, 0)
	assert.Equal(t, css.CssRuleList[0].Style.Styles["width"].Value, "100%")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["width"].Property, "width")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["width"].Important, 0)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "div .a")
}

func TestValuePx(t *testing.T) {
	css := Parse("div .a { font-size: 45px;}")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Value, "45px")
}

func TestValueEm(t *testing.T) {
	css := Parse("div .a { font-size: 45em;}")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["font-size"].Value, "45em")
}

func TestValueHex(t *testing.T) {
	css := Parse("div .a { color: #123456;}")
	assert.Equal(t, css.CssRuleList[0].Style.Styles["color"].Value, "#123456")
}

func TestValueRGBFunction(t *testing.T) {
	css := Parse(".color{ color: rgb(1,2,3);}")

	assert.Equal(t, css.CssRuleList[0].Style.Styles["color"].Value, "rgb(1,2,3)")
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, ".color")
}

func TestValueString(t *testing.T) {
	css := Parse("div .center { text-align: center; }")

	assert.Equal(t, css.CssRuleList[0].Style.Styles["text-align"].Value, "center")
}

func TestValueWhiteSpace(t *testing.T) {
	css := Parse(".div { padding: 10px 0 0 10px}")

	assert.Equal(t, "10px 0 0 10px", css.CssRuleList[0].Style.Styles["padding"].Value)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, ".div")
}

func TestValueMixed(t *testing.T) {
	css := Parse(`td {
			padding: 0 12px 0 10px;
    		border-right: 1px solid white
		}`)

	assert.Equal(t, "0 12px 0 10px", css.CssRuleList[0].Style.Styles["padding"].Value)
	assert.Equal(t, "1px solid white", css.CssRuleList[0].Style.Styles["border-right"].Value)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "td")
}

func TestQuoteValue(t *testing.T) {
	css := Parse(`blockquote {
    				font-family: "Source Sans Pro", Arial, sans-serif;
			    	font-size: 27px;
			    	line-height: 35px;}`)

	assert.Equal(t, "\"Source Sans Pro\", Arial, sans-serif", css.CssRuleList[0].Style.Styles["font-family"].Value)
	assert.Equal(t, "27px", css.CssRuleList[0].Style.Styles["font-size"].Value)
	assert.Equal(t, "35px", css.CssRuleList[0].Style.Styles["line-height"].Value)
	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "blockquote")
}

func TestDashClassname(t *testing.T) {
	css := Parse(`.content {
    				padding: 0px;
						}
						.content-wrap {
					  padding: 2px;
						}`)

	assert.Equal(t, ".content", css.CssRuleList[0].Style.SelectorText)
	assert.Equal(t, ".content-wrap", css.CssRuleList[1].Style.SelectorText)
	assert.Equal(t, "0px", css.CssRuleList[0].Style.Styles["padding"].Value)
	assert.Equal(t, "2px", css.CssRuleList[1].Style.Styles["padding"].Value)
}

func TestNotSupportedAtRule(t *testing.T) {
	rules := []string{
		`@keyframes mymove {
			    0%   {top: 0px;}
			    25%  {top: 200px;}
			    50%  {top: 100px;}
			    75%  {top: 200px;}
			    100% {top: 0px;}
			}`,
		`@-webkit-keyframes mymove {
			    0%   {top: 0px;}
			    25%  {top: 200px;}
			    50%  {top: 100px;}
			    75%  {top: 200px;}
			    100% {top: 0px;}
			} `,
		`@counter-style winners-list {
			  system: fixed;
			  symbols: url(gold-medal.svg) url(silver-medal.svg) url(bronze-medal.svg);
			  suffix: " ";
			}`,
		`@namespace url(http://www.w3.org/1999/xhtml);`,
		`@document url(http://www.w3.org/),
               url-prefix(http://www.w3.org/Style/),
               domain(mozilla.org),
               regexp("https:.*")
			{

			  body { color: purple; background: yellow; }
			}`,
	}
	for _, rule := range rules {
		assert.Panics(t, func() {
			Parse(rule)
		})
	}
}

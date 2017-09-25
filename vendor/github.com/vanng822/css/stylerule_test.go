package css

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStyleRuleText(t *testing.T) {
	sr := CSSStyleRule{}
	sr.SelectorText = ".box"
	sr.Styles = make(map[string]*CSSStyleDeclaration)
	sr.Styles["width"] = NewCSSStyleDeclaration("width", "10px", 0)
	sr.Styles["height"] = NewCSSStyleDeclaration("height", "100px", 0)
	
	assert.Equal(t, sr.Text(), ".box {\nheight: 100px;\nwidth: 10px\n}")
}
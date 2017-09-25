package css

import (
	"github.com/gorilla/css/scanner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCharsetDoubleQ(t *testing.T) {
	css := Parse(`@charset "UTF-8";`)

	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "\"UTF-8\"")
	assert.Equal(t, css.CssRuleList[0].Type, CHARSET_RULE)
}

func TestCharsetSingleQ(t *testing.T) {
	css := Parse(`@charset 'iso-8859-15';`)

	assert.Equal(t, css.CssRuleList[0].Style.SelectorText, "'iso-8859-15'")
	assert.Equal(t, css.CssRuleList[0].Type, CHARSET_RULE)
}

func TestCharsetIgnore(t *testing.T) {
	css := parseCharset(scanner.New(` 'iso-8859-15'`))

	assert.Nil(t, css)
}

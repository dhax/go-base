package css

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRuleTypeString(t *testing.T) {
	assert.Equal(t, STYLE_RULE.Text(), "")
	assert.Equal(t, CHARSET_RULE.Text(), "@charset")
	assert.Equal(t, IMPORT_RULE.Text(), "@import")
	assert.Equal(t, MEDIA_RULE.Text(), "@media")
	assert.Equal(t, FONT_FACE_RULE.Text(), "@font-face")
	assert.Equal(t, PAGE_RULE.Text(), "@page")
}


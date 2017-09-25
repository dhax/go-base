package css

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeclWithImportan(t *testing.T) {
	decl := NewCSSStyleDeclaration("width", "100%", 1)
	assert.Equal(t, decl.Text(), "width: 100% !important")
}

func TestDeclWithoutImportan(t *testing.T) {
	decl := NewCSSStyleDeclaration("width", "100%", 0)
	assert.Equal(t, decl.Text(), "width: 100%")
}
package css

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//"fmt"
)

func TestParseBlock(t *testing.T) {
	css := ParseBlock(`
    				font-family: "Source Sans Pro", Arial, sans-serif;
			    	font-size: 27px;
			    	line-height: 35px;`)

	assert.Equal(t, len(css), 3)
	assert.Equal(t, "35px", css["line-height"].Value)
}

func TestParseBlockOneLine(t *testing.T) {
	css := ParseBlock("font-family: \"Source Sans Pro\", Arial, sans-serif; font-size: 27px;")

	assert.Equal(t, len(css), 2)
	assert.Equal(t, "27px", css["font-size"].Value)
	assert.Equal(t, "\"Source Sans Pro\", Arial, sans-serif", css["font-family"].Value)
}

func TestParseBlockBlankEnd(t *testing.T) {
	css := ParseBlock("font-size: 27px; width: 10px")
	
	assert.Equal(t, len(css), 2)
	assert.Equal(t, "27px", css["font-size"].Value)
	assert.Equal(t, "10px", css["width"].Value)
}

func TestParseBlockInportant(t *testing.T) {
	css := ParseBlock("font-size: 27px; width: 10px !important")
	
	assert.Equal(t, len(css), 2)
	assert.Equal(t, "27px", css["font-size"].Value)
	assert.Equal(t, "10px", css["width"].Value)
	assert.Equal(t, 1, css["width"].Important)
}

func TestParseBlockWithBraces(t *testing.T) {
	css := ParseBlock("{ font-size: 27px; width: 10px }")
	
	assert.Equal(t, len(css), 2)
	assert.Equal(t, "27px", css["font-size"].Value)
	assert.Equal(t, "10px", css["width"].Value)
}
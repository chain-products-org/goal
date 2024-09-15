package fontx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	path := "../testdata/fonts/Encode_Sans/EncodeSans-Bold.ttf"
	ft, err := Load(path)
	assert.True(t, err == nil)
	assert.True(t, ft != nil)
}

func TestFace(t *testing.T) {
	path := "../testdata/fonts/Encode_Sans/EncodeSans-Bold.ttf"
	ft, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	fa := Face(ft)
	assert.True(t, fa != nil)
}

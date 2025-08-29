package imagex_test

import (
	"bytes"
	"image/jpeg"
	"os"
	"testing"

	"github.com/chain-products-org/goal/errorx"
	"github.com/chain-products-org/goal/imagex"
)

var (
	pic  = "testdata/image/blue-purple-pink.png"
	dpic = "testdata/image/blue-purple-pink_thumb.png"
)

func TestResize(t *testing.T) {
	f, err := os.Open(pic)
	errorx.Throw(err)
	img, err := jpeg.Decode(f)
	errorx.Throw(err)
	dimg := imagex.Resize(200, 200, img)

	errorx.Throw(err)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, dimg, &jpeg.Options{Quality: 100})
	err = os.WriteFile(pic, buf.Bytes(), os.ModePerm)
	errorx.Throw(err)
}

func TestThumbnail(t *testing.T) {
	f, err := os.Open(pic)
	errorx.Throw(err)
	img, err := jpeg.Decode(f)
	errorx.Throw(err)
	dimg := imagex.Thumbnail(400, 400, img)

	errorx.Throw(err)
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, dimg, &jpeg.Options{Quality: 100})
	err = os.WriteFile(dpic, buf.Bytes(), os.ModePerm)
	errorx.Throw(err)
}

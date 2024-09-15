package imagex_test

import (
	"fmt"
	"image"
	"image/color"
	"testing"

	"github.com/gophero/goal/fontx"
	"github.com/gophero/goal/imagex"
	"github.com/stretchr/testify/assert"
)

func TestDrawMonochrome(t *testing.T) {
	w, h := 100, 200
	color := imagex.NewRGBA(0, 0, 255, 255)
	img := imagex.DrawPureRect(w, h, color)
	bounds := img.Bounds()
	assert.NotNil(t, img, "img should not be nil")
	assert.True(t, img.At(bounds.Min.X, bounds.Min.Y) == color)
	assert.True(t, img.At(0, h/2) == color)
	assert.True(t, img.At(w/2, h/2) == color)
	assert.True(t, img.At(bounds.Min.X, bounds.Min.Y) == color)
	assert.True(t, img.At(bounds.Max.X-1, bounds.Max.Y-1) == color) // 0 - w - 1
	assert.True(t, img.Bounds().Dy() == h)
	assert.True(t, img.Bounds().Dx() == w)
}

func TestDrawPureCircle(t *testing.T) {
	r := 50
	tp := imagex.NewRGBA(0, 0, 0, 0)
	clr := imagex.NewRGBA(255, 0, 0, 255)
	img := imagex.DrawPureCircle(clr, r)
	assert.True(t, img != nil)
	bounds := img.Bounds()
	imagex.Save(img, "t.png")
	fmt.Println(img.At(bounds.Min.X, bounds.Min.Y))
	fmt.Println(color.Transparent)
	assert.True(t, img.At(bounds.Min.X, bounds.Min.Y) == tp)
	assert.True(t, img.At(0, r) == clr)
	assert.True(t, img.At(r, 0) == clr)
	assert.True(t, img.At(r, r/2) == clr)
	assert.True(t, img.At(r/2, r) == clr)
	assert.True(t, img.At(bounds.Max.X-1, bounds.Max.Y-1) == tp) // 0 - w - 1
	assert.True(t, img.Bounds().Dy() == r*2)
	assert.True(t, img.Bounds().Dx() == r*2)
}

func TestCropCircle(t *testing.T) {
	blue := imagex.NewRGBA(0, 0, 255, 255)
	// red := image.NewRGBA(image.Rect(255, 0, 0, 255))
	src := imagex.DrawPureRect(100, 200, blue)
	dest := imagex.CropCircle(src)
	// imagex.Save(dest, "crop_circle.png")
	assert.True(t, dest != nil)
	assert.True(t, dest.Bounds().Dx() > 0 && dest.Bounds().Dy() > 0)
}

func TestMontage(t *testing.T) {
	imgs := []image.Image{
		imagex.DrawPureRect(50, 10, imagex.NewRGBA(255, 0, 0, 255)),
		imagex.DrawPureRect(50, 25, imagex.NewRGBA(255, 255, 0, 255)),
		imagex.DrawPureRect(50, 50, imagex.NewRGBA(255, 255, 255, 255)),
		imagex.DrawPureRect(10, 50, imagex.NewRGBA(255, 0, 255, 255)),
		imagex.DrawPureRect(25, 50, imagex.NewRGBA(0, 255, 255, 255)),
		imagex.DrawPureRect(50, 50, imagex.NewRGBA(0, 0, 255, 255)),
		imagex.DrawPureRect(100, 100, imagex.NewRGBA(0, 255, 0, 255)),
		imagex.DrawPureRect(100, 50, imagex.NewRGBA(0, 0, 0, 255)),
	}
	img := imagex.Montage(imgs)
	// imagex.Save(img, "default.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)
	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod())
	// imagex.Save(img, "col.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)
	img = imagex.Montage(imgs, imagex.MontageOpts.Len(1))
	// imagex.Save(img, "col1.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)
	img = imagex.Montage(imgs, imagex.MontageOpts.Len(2))
	// imagex.Save(img, "col2.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)
	img = imagex.Montage(imgs, imagex.MontageOpts.Len(3))
	// imagex.Save(img, "col3.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(4))
	// imagex.Save(img, "col4.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(5))
	// imagex.Save(img, "col5.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(6))
	// imagex.Save(img, "col6.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(7))
	// imagex.Save(img, "col7.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod())
	// imagex.Save(img, "row.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(1))
	// imagex.Save(img, "row1.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(2))
	// imagex.Save(img, "row2.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(3))
	// imagex.Save(img, "row3.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(4))
	// imagex.Save(img, "row4.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(5))
	// imagex.Save(img, "row5.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(6))
	// imagex.Save(img, "row6.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(7))
	// imagex.Save(img, "row7.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Compress())
	// imagex.Save(img, "colc.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(1), imagex.MontageOpts.Compress())
	// imagex.Save(img, "col1c.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(3), imagex.MontageOpts.Compress())
	// imagex.Save(img, "col3.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(4), imagex.MontageOpts.Compress())
	// imagex.Save(img, "col4.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.Len(7), imagex.MontageOpts.Compress())
	// imagex.Save(img, "col7c.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Compress())
	// imagex.Save(img, "rowc.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(1), imagex.MontageOpts.Compress())
	// imagex.Save(img, "row1c.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(3), imagex.MontageOpts.Compress())
	// imagex.Save(img, "row3c.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(4), imagex.MontageOpts.Compress())
	// imagex.Save(img, "row4c.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)

	img = imagex.Montage(imgs, imagex.MontageOpts.RowMod(), imagex.MontageOpts.Len(7), imagex.MontageOpts.Compress())
	// imagex.Save(img, "row7c.png")
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dx() > 0 && img.Bounds().Dy() > 0)
}

func TestDrawText(t *testing.T) {
	text := "hello"
	w, h, color := 100, 200, imagex.NewRGBA(255, 255, 0, 255)
	src := imagex.DrawPureRect(w, h, color)
	ft, err := fontx.Load("../testdata/fonts/Encode_Sans/EncodeSans-Bold.ttf")
	if err != nil {
		t.Fatalf("load font error: %v", err)
	}
	fontColor := imagex.NewRGBA(0, 255, 0, 255)
	dest := imagex.DrawText(src, text, fontColor, fontx.Face(ft, fontx.Size(20)), 0, 0)
	imagex.Save(dest, "text.png")
}

package imagex_test

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gophero/goal/imagex"
	"github.com/stretchr/testify/assert"
)

func newBlueImg() image.Image {
	w, h := 100, 200
	color := imagex.NewRGBA(0, 0, 255, 255)
	img := imagex.DrawPureRect(w, h, color)
	return img
}

func TestSave(t *testing.T) {
	save("t.jpg", nil, t)
	save("t.png", nil, t)
	save("t.gif", nil, t)
	save("t.bmp", nil, t)
}

func save(name string, opt imagex.Option, t *testing.T) {
	img := newBlueImg()
	err := imagex.Save(img, name, opt)
	if err != nil {
		t.Errorf("save test failed: %v", err)
		t.FailNow()
	}
	f, err := os.Open(name)
	if err != nil {
		t.Errorf("save test failed, file can not open: %v", name)
	}
	assert.True(t, f != nil)
	f.Close()
	os.Remove(name)
}

func TestSaveWithOpt(t *testing.T) {
	save("t.jpg", imagex.NewJpgOption(20), t)
	save("t.png", nil, t)
	save("t.gif", imagex.NewGifOption(256, nil, nil), t)
	save("t.bmp", nil, t)
}

func TestDecodeFile(t *testing.T) {
	p := "../testdata/image/blue-purple-pink.png"
	f, err := os.Open(p)
	if err != nil {
		t.Errorf("open file %s failed: %v", p, err)
		t.FailNow()
	}
	img, err := imagex.Decode(f)
	if err != nil {
		t.Errorf("decode image file %s error: %v", p, err)
		t.FailNow()
	}
	assert.True(t, img != nil)
}

func TestDecodeFromUrl(t *testing.T) {
	// NOTE: raw image saved in github is from raw.githubusercontent.com, NOT https://github.com/gophero/goal/blob/dev/testdata/image/blue-purple-pink.png
	// which will report an error: image is not an png/unkown format
	url := "https://raw.githubusercontent.com/gophero/goal/dev/testdata/image/blue-purple-pink.png"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("request image url %s error: %v", url, err)
		t.FailNow()
	}
	assert.True(t, resp.Body != nil)
	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read body failed: %v", err)
		t.FailNow()
	}
	resp.Body.Close()
	r := bytes.NewReader(bs)
	img, err := imagex.Decode(r)
	if err != nil {
		t.Errorf("decode image from url %s error: %v", url, err)
		t.FailNow()
	}
	assert.True(t, img != nil)
	assert.True(t, img.Bounds().Dy() > 0 && img.Bounds().Dx() > 0)
}

func TestDecodeConfigFromFile(t *testing.T) {
	f, err := os.Open("../testdata/image/blue-purple-pink.png")
	if err != nil {
		t.Errorf("open file failed, error: %v", err)
		t.FailNow()
	}

	if err != nil {
		t.Errorf("read body error: %v", err)
		t.FailNow()
	}
	c, ft, err := image.DecodeConfig(f)
	if err != nil {
		t.Errorf("decode image config failed, error: %v", err)
		t.FailNow()
	}
	fmt.Println("config:", c.ColorModel, c.Width, c.Height)
	fmt.Println("format:", ft)
	assert.True(t, ft == "png")

	f.Seek(0, 0)
	img, err := png.Decode(f)
	if err != nil {
		t.Errorf("decode png failed: %v", err)
		t.FailNow()
	}
	assert.True(t, img != nil)
}

func TestDecodeConfigFromUrl(t *testing.T) {
	url := "https://raw.githubusercontent.com/gophero/goal/dev/testdata/image/blue-purple-pink.png"
	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("request image url %s error: %v", url, err)
		t.FailNow()
	}

	bs, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("read body error: %v", err)
		t.FailNow()
	}
	r := bytes.NewReader(bs)
	r.Seek(0, 0)
	c, ft, err := image.DecodeConfig(r)
	if err != nil {
		t.Errorf("decode image config failed, error: %v", err)
		t.FailNow()
	}
	fmt.Println("config:", c.ColorModel, c.Width, c.Height)
	fmt.Println("format:", ft)
	assert.True(t, ft == "png")

	r.Seek(0, 0)
	img, err := png.Decode(r)
	if err != nil {
		t.Errorf("decode png failed: %v", err)
		t.FailNow()
	}
	assert.True(t, img != nil)

	// imagex.Save(img, "haha.png")
}

func TestDecodeFromUrl1(t *testing.T) {
	url := "https://raw.githubusercontent.com/gophero/goal/dev/testdata/image/blue-purple-pink.png"
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	m, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	g := m.Bounds()

	// Get height and width
	height := g.Dy()
	width := g.Dx()

	// The resolution is height x width
	resolution := height * width

	// Print results
	fmt.Println(resolution, "pixels")
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}

func TestDrawCircleLine(t *testing.T) {
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	drawCircle(img, 40, 40, 30, color.RGBA{255, 0, 0, 255})

	buf := &bytes.Buffer{}
	if err := png.Encode(buf, img); err != nil {
		panic(err)
	}
	if err := os.WriteFile("circle.png", buf.Bytes(), 0666); err != nil {
		panic(err)
	}
}

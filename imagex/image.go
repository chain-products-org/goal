package imagex

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/chain-products-org/goal/valuex"
	"golang.org/x/image/bmp"
	"golang.org/x/image/draw"
	"golang.org/x/image/webp"
)

const (
	JPG  = "jpg"
	JPEG = "jpeg"
	PNG  = "png"
	GIF  = "gif"
	BMP  = "bmp"
	WEBP = "webp"
)

var UnsupportExtError = fmt.Errorf("unsupported image extension")

type Option interface{}

func NewJpgOption(quality int) Option {
	return &jpeg.Options{Quality: quality}
}

func DefaultQuality() Option {
	return &jpeg.Options{Quality: jpeg.DefaultQuality}
}

func HighQulity() Option {
	return &jpeg.Options{Quality: 100}
}

func LowQulity() Option {
	return &jpeg.Options{Quality: 1}
}

func NewGifOption(numColors int, quantizer draw.Quantizer, drawer draw.Drawer) Option {
	return &gif.Options{NumColors: numColors, Quantizer: quantizer, Drawer: drawer}
}

func Save(img image.Image, path string, options ...Option) error {
	ext := FmtExt(path)
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	var option Option
	if len(options) > 0 {
		option = options[0]
	}
	switch ext {
	case JPG, JPEG:
		if option == nil {
			option = DefaultQuality()
		}
		err = jpeg.Encode(f, img, option.(*jpeg.Options))
	case PNG:
		err = png.Encode(f, img)
	case GIF:
		if option != nil {
			err = gif.Encode(f, img, valuex.Def(option == nil, nil, option.(*gif.Options)))
		} else {
			err = gif.Encode(f, img, nil)
		}
	case BMP:
		err = bmp.Encode(f, img)
	case WEBP:
		// TODO: not supported yet
		fallthrough
	default:
		return UnsupportExtError
	}
	if err != nil {
		return err
	}
	return f.Close()
}

func FmtExt(path string) string {
	ext := filepath.Ext(path)
	path = strings.TrimRight(path, ext)
	ext = strings.ToLower(ext)
	return strings.TrimLeft(ext, ".")
}

func Decode(r io.Reader) (image.Image, error) {
	// make r to read twice to decode config and decode image
	bs, err := io.ReadAll(r)
	br := bytes.NewReader(bs)
	_, ft, err := image.DecodeConfig(br)
	if err != nil {
		return nil, err
	}
	// can replace with image.Decode but need to import png,gif,jpeg,bmp package
	// so i choose to use the implements package to decode
	r = bytes.NewReader(bs)
	switch ft {
	case JPG, JPEG:
		return jpeg.Decode(r)
	case PNG:
		return png.Decode(r)
	case GIF:
		return gif.Decode(r)
	case BMP:
		return bmp.Decode(r)
	case WEBP:
		return webp.Decode(r)
	}
	return nil, nil
}

func DecodeFile(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return Decode(f)
}

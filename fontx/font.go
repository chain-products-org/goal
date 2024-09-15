package fontx

import (
	"io"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func Load(path string) (*truetype.Font, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return truetype.Parse(bs)
}

type FaceOption func(opts *truetype.Options)

func Size(size float64) FaceOption {
	return func(opts *truetype.Options) {
		opts.Size = size
	}
}

func DPI(dpi float64) FaceOption {
	return func(opts *truetype.Options) {
		opts.DPI = dpi
	}
}

func Hinting(hinting font.Hinting) FaceOption {
	return func(opts *truetype.Options) {
		opts.Hinting = hinting
	}
}

func HintingVertical() FaceOption {
	return func(opts *truetype.Options) {
		opts.Hinting = font.HintingVertical
	}
}

func HintingFull() FaceOption {
	return func(opts *truetype.Options) {
		opts.Hinting = font.HintingFull
	}
}

func GlyphCacheEntries(g int) FaceOption {
	return func(opts *truetype.Options) {
		opts.GlyphCacheEntries = g
	}
}

func SubPixelsX(x int) FaceOption {
	return func(opts *truetype.Options) {
		opts.SubPixelsX = x
	}
}

func SubPixelsY(y int) FaceOption {
	return func(opts *truetype.Options) {
		opts.SubPixelsY = y
	}
}

func Face(f *truetype.Font, opts ...FaceOption) font.Face {
	fopt := &truetype.Options{}
	for _, opt := range opts {
		opt(fopt)
	}
	return truetype.NewFace(f, fopt)
}

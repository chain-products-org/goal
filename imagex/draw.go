package imagex

import (
	"image"
	"image/color"
	"image/draw"

	"github.com/gophero/goal/mathx"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type Direction int

const (
	directionColumn Direction = iota
	directionRow
)

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if xx*xx+yy*yy < rr*rr {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

// DrawPureRect draw a pure color rectangle image using the given size and color.
func DrawPureRect(w, h int, color color.RGBA) image.Image {
	m := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.Draw(m, m.Bounds(), &image.Uniform{color}, image.Point{}, draw.Src)
	return m
}

func DrawPureCircle(color color.RGBA, radius int) image.Image {
	rect := image.Rect(0, 0, radius*2, radius*2)
	dest := image.NewRGBA(rect) // a transpanret image
	circleImg := image.NewUniform(color)
	draw.DrawMask(dest, dest.Bounds(), circleImg, image.Point{}, &circle{p: image.Point{radius, radius}, r: radius}, image.Point{}, draw.Over)
	return dest
}

func CropCircle(img image.Image) image.Image {
	center := image.Point{img.Bounds().Dx() / 2, img.Bounds().Dy() / 2}
	radius := mathx.Minn(img.Bounds().Dx(), img.Bounds().Dy()) / 2
	rect := image.Rect(center.X-radius, center.Y-radius, center.X+radius, center.Y+radius)
	dest := image.NewRGBA(rect)
	draw.DrawMask(dest, dest.Bounds(), img, image.Point{}, &circle{p: image.Point{rect.Dx() / 2, rect.Dy() / 2}, r: radius}, image.Point{}, draw.Over)
	return dest
}

var MontageOpts = montageOptions{}

type (
	MontageOption  func(opts *montageOptions)
	montageOptions struct {
		d        Direction
		len      int
		compress bool
	}
)

func (montageOptions) RowMod() MontageOption {
	return func(opts *montageOptions) {
		opts.d = directionRow
	}
}

func (montageOptions) ColMod() MontageOption {
	return func(opts *montageOptions) {
		opts.d = directionColumn
	}
}

func (montageOptions) Len(n int) MontageOption {
	return func(opts *montageOptions) {
		opts.len = n
	}
}

func (montageOptions) Compress() MontageOption {
	return func(opts *montageOptions) {
		opts.compress = true
	}
}

func Montage(imgs []image.Image, opts ...MontageOption) image.Image {
	opt := new(montageOptions)
	for _, v := range opts {
		v(opt)
	}
	if opt.len == 0 {
		opt.len = len(imgs)
	}

	rowMod := opt.d == directionRow
	if opt.compress {
		var mw, mh int = imgs[0].Bounds().Dx(), imgs[0].Bounds().Dy()
		for _, img := range imgs {
			mw = mathx.Minn(mw, img.Bounds().Dx())
			mh = mathx.Minn(mh, img.Bounds().Dy())
		}
		newimgs := []image.Image{}
		for _, img := range imgs {
			w, h := img.Bounds().Dx(), img.Bounds().Dy()
			if rowMod {
				img = Thumbnail(uint(w), uint(mh), img)
			} else {
				img = Thumbnail(uint(mw), uint(h), img)
			}
			newimgs = append(newimgs, img)
		}
		imgs = newimgs
	}

	unit := len(imgs) / opt.len
	if len(imgs)%opt.len > 0 {
		unit += 1
	}
	var tws, ths [][]int
	x := 0
	for i, img := range imgs {
		if i == 0 {
			tws = append(tws, []int{})
			ths = append(ths, []int{})
		} else {
			if i%opt.len == 0 {
				x += 1
				tws = append(tws, []int{})
				ths = append(ths, []int{})
			}
		}
		tws[x] = append(tws[x], img.Bounds().Dx())
		ths[x] = append(ths[x], img.Bounds().Dy())
	}
	var tw, th int
	var maxUnit []int
	if rowMod {
		for _, vs := range tws {
			tw = mathx.Maxn(tw, mathx.Sumn(vs...))
		}
		for _, vs := range ths {
			th += mathx.Maxn(vs...)
			maxUnit = append(maxUnit, mathx.Maxn(vs...))
		}
	} else {
		for _, vs := range tws {
			tw += mathx.Maxn(vs...)
			maxUnit = append(maxUnit, mathx.Maxn(vs...))
		}
		for _, vs := range ths {
			th = mathx.Maxn(th, mathx.Sumn(vs...))
		}
	}
	dest := image.NewRGBA(image.Rect(0, 0, tw, th))
	offx, offy := 0, 0
	for i, img := range imgs {
		draw.Draw(dest, img.Bounds().Add(image.Point{offx, offy}), img, image.Point{}, draw.Over)
		idx := i / opt.len
		if rowMod {
			if (i+1)%opt.len == 0 {
				offx = 0
				offy += maxUnit[idx]
			} else {
				offx += img.Bounds().Dx()
			}
		} else {
			if (i+1)%opt.len == 0 && i+1 < len(imgs) {
				offy = 0
				offx += maxUnit[idx]
			} else {
				offy += img.Bounds().Dy()
			}
		}
	}
	return dest
}

func DrawText(src image.Image, text string, color color.Color, face font.Face, x, y int) image.Image {
	// new RGBA which implements draw.Image
	dest := image.NewRGBA(src.Bounds())
	draw.Draw(dest, dest.Bounds(), src, image.Point{}, draw.Src)
	drawer := &font.Drawer{Dst: dest, Src: image.NewUniform(color), Face: face}
	drawer.Dot = fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
	drawer.DrawString(text)
	return dest
}

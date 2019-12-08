package image

import (
	"fmt"
	"image"
	"image/color"
)

type Color int

const (
	Black       Color = 0
	White       Color = 1
	Transparent Color = 2
)

func (c Color) RGBA() (r, g, b, a uint32) {
	switch c {
	case Black:
		return 0, 0, 0, 0xffff
	case White:
		return 0xffff, 0xffff, 0xffff, 0xffff
	case Transparent:
		return 0, 0, 0, 0
	default:
		panic("Invalid color")
	}
}

type Layer [][]Color

func (l Layer) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color {
		r, g, b, a := c.RGBA()
		if a < 0x8000 {
			return Transparent
		}

		i := (r + g + b) / 3
		if i < 0x8000 {
			return Black
		}

		return White
	})
}

func (l Layer) Bounds() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{len(l[0]), len(l)},
	}
}

func (l Layer) At(x, y int) color.Color {
	return l[y][x]
}

type Image []Layer

func decodeRow(s string) ([]Color, error) {
	row := make([]Color, len(s))

	for i, ch := range s {
		if ch < '0' || ch > '9' {
			return nil, fmt.Errorf("Invalid character %q", string(ch))
		}

		row[i] = Color(int(ch) - '0')
	}

	return row, nil
}

func decodeLayer(s string, width int) (Layer, error) {
	layer := make(Layer, 0, len(s)/width)

	for i := 0; i < len(s); i += width {
		row, err := decodeRow(s[i : i+width])
		if err != nil {
			return nil, err
		}
		layer = append(layer, row)
	}

	return layer, nil
}

func Decode(s string, width, height int) (Image, error) {
	size := width * height

	if len(s)%size != 0 {
		return nil, fmt.Errorf("Invalid number of characters")
	}

	img := make(Image, 0, len(s)/size)

	for i := 0; i < len(s); i += size {
		layer, err := decodeLayer(s[i:i+size], width)
		if err != nil {
			return nil, err
		}
		img = append(img, layer)
	}

	return img, nil
}

func (i Image) Flatten() Layer {
	height := len(i[0])
	width := len(i[0][0])

	res := make(Layer, height)
	for i := range res {
		res[i] = make([]Color, width)
		for j := range res[i] {
			res[i][j] = Transparent
		}
	}

	for _, layer := range i {
		for y, row := range layer {
			for x, px := range row {
				if res[y][x] == Transparent {
					res[y][x] = px
				}
			}
		}
	}

	return res
}

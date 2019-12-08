package image

import (
	"fmt"
)

type Image [][][]int

func decodeRow(s string) ([]int, error) {
	row := make([]int, len(s))

	for i, ch := range s {
		if ch < '0' || ch > '9' {
			return nil, fmt.Errorf("Invalid character %q", string(ch))
		}

		row[i] = int(ch) - '0'
	}

	return row, nil
}

func decodeLayer(s string, width int) ([][]int, error) {
	layer := make([][]int, 0, len(s)/width)

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

	data := make([][][]int, 0, len(s)/size)

	for i := 0; i < len(s); i += size {
		layer, err := decodeLayer(s[i:i+size], width)
		if err != nil {
			return nil, err
		}
		data = append(data, layer)
	}

	return Image(data), nil
}

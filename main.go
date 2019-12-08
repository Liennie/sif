package main

import (
	"fmt"
	"io/ioutil"

	"github.com/liennie/sif/image"
)

func main() {
	data, err := ioutil.ReadFile("data.txt")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	img, err := image.Decode(string(data), 25, 6)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	min := 25 * 6
	product := 0

	for _, layer := range img {
		counts := make([]int, 3)
		for _, row := range layer {
			for _, px := range row {
				if px < 0 || px > 2 {
					fmt.Println("Invalid character")
					return
				}

				counts[px]++
			}
		}

		if counts[0] < min {
			min = counts[0]
			product = counts[1] * counts[2]
		}
	}

	fmt.Println(product)
}

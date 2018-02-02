package main

import "golang.org/x/tour/pic"
import "math/rand"

func Pic(dx, dy int) [][]uint8 {
	var res [][]uint8
	for i := 0; i < dy; i++ {
		var t []uint8
		for j := 0; j < dx; j++ {
			if i == j {
				t = append(t, uint8(rand.Intn(256)))
			} else {
				t = append(t, 255)
			}
		}
		res = append(res, t)
	}
	return res
}

func main() {
	pic.Show(Pic)
}

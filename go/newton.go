package main

import (
	"fmt"
	"math"
)

func main() {
	fmt.Println(Sqrt(2))
}

func Sqrt(x float64) float64 {
	var z float64
	z = x
	for i := 1; i < 10; i++ {
		z = z - (math.Pow(z, 2)-x)/(2*z)
		fmt.Println(z)
	}
	return z
}

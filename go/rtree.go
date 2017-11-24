package main

import (
	"github.com/dhconnelly/rtreego"
	"fmt"
)

func main() {
	rt := rtreego.NewTree(2, 25, 50)
	p1 := rtreego.Point{0.4, 0.5}
	p2 := rtreego.Point{6.2, -3.4}

	r1, _ := rtreego.NewRect(p1, []float64{1, 2})
	r2, _ := rtreego.NewRect(p2, []float64{1.7, 2.7})

	rt.Insert(&Thing{r1, "foo"})
	rt.Insert(&Thing{r2, "bar"})

	size := rt.Size()

	fmt.Println(size)
	q := rtreego.Point{6.5, -2.47}
	k := 5
	results := rt.NearestNeighbors(k, q)
	fmt.Println(results)
}

type Thing struct {
	where *rtreego.Rect
	name string
}

func (t *Thing) Bounds() *rtreego.Rect {
	return t.where
}

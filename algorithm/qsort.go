package main

import (
    "fmt"
    "math/rand"
    "time"
)

func swap(a int, b int) (int, int) {
    return b, a
}

func partition(aris []int, begin int, end int) int {
    pvalue := aris[begin]
    i := begin
    j := begin + 1
    for j < end {
        if aris[j] < pvalue {
            i++
            aris[i], aris[j] = swap(aris[i], aris[j])
        }
        j++
    }
    aris[i], aris[begin] = swap(aris[i], aris[begin])
    return i
}

func quickSort(aris []int, begin int, end int) {
    if begin+1 < end {
        mid := partition(aris, begin, end)
        quickSort(aris, begin, mid)
        quickSort(aris, mid+1, end)
    }
}

func randArray(aris []int) {
    l := len(aris)
    for i := 0; i < l; i++ {
        r := rand.New(rand.NewSource(time.Now().UnixNano()))
        aris[i] = r.Intn(1000)
    }
}

func main() {
    intas := make([]int, 10)
    randArray(intas)
    fmt.Println(intas)

    quickSort(intas, 0, 10)
    fmt.Println(intas)
}

package main

import (
	"fmt"
	"github.com/spf13/cast"
)

type Animal interface {
	Noise() string
}

type Dog struct{}

func (Dog) Noise() string {
	return "Woof!"
}

func PrintNoises(as []Animal) {
	for _, a := range as {
		fmt.Println(a.Noise())
	}
}
func main() {
	println(cast.ToString("mayonegg"))

	animals := []Animal{Dog{}}
	PrintNoises(animals)

	/*
	dogs := []Dog{Dog{}}
	PrintNoises(dogs)

	 */
	dogs := []Dog{Dog{}}
	// 新逻辑: 把 dogs 的切片转换成 animals 的切片
	animals = []Animal{}
	for _, d := range dogs {
		animals = append(animals, Animal(d))
	}
	PrintNoises(animals)
}


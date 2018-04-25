package main

import (
	"bytes"
	"fmt"
	"encoding/gob"
	"log"
	"os"
)

type P struct {
	X, Y, Z int
	Name    string
}

type Q struct {
	Z, X, Y *int32
	Name    string
}

func main1() {
	// Initialize the encoder and decoder.  Normally enc and dec would be
	// bound to network connections and the encoder and decoder would
	// run in different processes.
	var network bytes.Buffer        // Stand-in for a network connection
	enc := gob.NewEncoder(&network) // Will write to network.
	dec := gob.NewDecoder(&network) // Will read from network.
	// Encode (send) the value.
	err := enc.Encode(P{3, 4, 5, "Pythagoras"})
	if err != nil {
		log.Fatal("encode error:", err)
	}
	// Decode (receive) the value.
	var q Q
	err = dec.Decode(&q)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Printf("%q: {%d,%d,%d}\n", q.Name, *q.X, *q.Y, *q.Z)
}

// Output:   "Pythagoras": {3,4}

type Address struct {
	Type    string
	City    string
	Country string
}

func (a *Address) sayHello() {
	fmt.Println("I am in " + a.City)
}

type VCard struct {
	FirstName string
	LastName  string
	Addresses []*Address
	Remark    string
}

var content string

func main2() {
	pa := &Address{"private", "Aartselaar", "Belgium"}
	wa := &Address{"work", "Boom", "Belgium"}
	vc := VCard{"Jan", "Kersschot", []*Address{pa, wa}, "none"}
	// fmt.Printf("%v: \n", vc) // {Jan Kersschot [0x126d2b80 0x126d2be0] none}:
	// using an encoder:
	file, _ := os.OpenFile("vcard.gob", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := gob.NewEncoder(file)
	err := enc.Encode(vc)
	if err != nil {
		log.Println("Error in encoding gob")
	}
}

func main() {
	file, _ := os.OpenFile("vcard.gob", os.O_RDONLY, 0666)
	dec := gob.NewDecoder(file)
	var vc VCard
	err := dec.Decode(&vc)
	if err != nil {
		log.Fatal("decode error:", err)
	}
	fmt.Println(vc)
	fmt.Println(vc.Addresses[0].City)
	vc.Addresses[0].sayHello()
}

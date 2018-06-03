package main

import (
	"fmt"
	"log"

	"github.com/valyala/fastjson"
)

func main()  {
	var p fastjson.Parser
	v, err := p.Parse(`{"foo":"bar", "baz": 123}`)
	if err != nil {
		log.Fatalf("cannot parse json: %s", err)
	}

	fmt.Printf("foo=%s, baz=%d\n", v.GetStringBytes("foo"), v.GetInt("baz"))

	s := `{
		"obj": { "foo": 1234 },
		"arr": [ 23,4, "bar" ],
		"str": "foobar"
	}`

	v, err = p.Parse(s)
	if err != nil {
		log.Fatalf("cannot parse json: %s", err)
	}
	o, err := v.Object()
	if err != nil {
		log.Fatalf("cannot obtain object from json value: %s", err)
	}

	o.Visit(func(k []byte, v *fastjson.Value) {
		switch string(k) {
		case "obj":
			fmt.Printf("object %s\n", v)
		case "arr":
			fmt.Printf("array %s\n", v)
		case "str":
			fmt.Printf("string %s\n", v)
		}
	})
}

// refer https://godoc.org/github.com/valyala/fastjson
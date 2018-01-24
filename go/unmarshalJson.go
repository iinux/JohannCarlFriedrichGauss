package main

import (
	"time"
	"fmt"
	"encoding/json"
)

type Time struct {
	time.Time
}

// returns time.Now() no matter what!
func (t *Time) UnmarshalJSON(b []byte) error {
	// you can now parse b as thoroughly as you want

	*t = Time{time.Now()}
	return nil
}

type Config struct {
	T Time
}

func main() {
	c := Config{}

	json.Unmarshal([]byte(`{"T": "bad-time"}`), &t c)

	fmt.Printf("%+v\n", c)
}

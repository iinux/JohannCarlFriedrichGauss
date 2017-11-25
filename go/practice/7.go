package main

type Resource struct {
	id int
}

type User struct {
	*Resource
	// Resource // Error: duplicate field Resource
	name string
}

func main() {

	u := User{
		&Resource{1},
		"Administrator",
	}

	println(u.id)
	println(u.Resource.id)
}

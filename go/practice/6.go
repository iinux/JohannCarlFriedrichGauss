package main

type Resource struct {
	id   int
	name string
}

type Classify struct {
	id int
}

type User struct {
	Resource // Resource.id 与 Classify.id 处于同一层次。
	Classify
	name string // 遮蔽 Resource.name。
}

func main() {

	u := User{
		Resource{1, "people"},
		Classify{100},
		"Jack",
	}

	println(u.name)          // User.name: Jack
	println(u.Resource.name) // people

	// println(u.id) // Error: ambiguous selector u.id
	println(u.Classify.id) // 100
}

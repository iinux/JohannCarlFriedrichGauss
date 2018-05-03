package main

import (
	"strings"
	"fmt"
)

func main() {
	judge([]string{
		"/*Test program */",
		"int main()",
		"{ ",
		"  // variable declaration ",
		"int a, b, c;",
		"/* This is a test",
		"   multiline  ",
		"   comment for ",
		"   testing */",
		"a = b + c;",
		"}"})
}

func judge(source []string)  {
	res := removeComments(source)
	for _, s := range res {
		fmt.Println(s)
	}
}

func removeComments(source []string) []string {
	var res []string
	for _,s := range source {
		i := strings.Index(s, "//")
		if i >= 0{
			s = s[0:i]
		}
		if s == "" {
			continue
		}
		res = append(res, s)
	}

	return res
}

package main

import (
	"strings"
	"fmt"
)

func main() {
	judge([]string{
		"a/*/Test program */b",
		"int main()",
		"{ ",
		"  // variable declaration ",
		"int a, b, c;",
		"a/* This is a test",
		"   multiline  ",
		"   comment for ",
		"   testing */b",
		"a = b + c;",
		"}"})
	judge([]string{
		"main() {",
		"/* here is commments",
		"  // still comments */",
		"   double s = 33;",
		"   cout << s;",
		"}",
	})
	judge([]string{
		"void func(int k) {",
		"// this function does nothing /*",
		"   k = k*2/4;",
		"   k = k/2;*/",
		"}",
	})
	judge([]string{
		"main() {",
		"  Node* p;",
		"  /* declare a Node",
		"  /*float f = 2.0",
		"   p->val = f;",
		"   /**/",
		"   p->val = 1;",
		"   //*/ cout << success;*/",
		"}",
		" ",
	}) // ["main() {","  Node* p;","  ","   p->val = 1;","   ","}"," "]

}

func judge(source []string) {
	res := removeComments(source)
	for _, s := range res {
		fmt.Println(s)
	}
	fmt.Println("===")
}

func removeComments(source []string) []string {
	var res []string
	var line string
	var inComment bool
	for _, s := range source {
		i := strings.Index(s, "//")
		if i >= 0 && !inComment{
			s = s[0:i]
		}

		i = strings.Index(s, "/*")
		j := strings.Index(s, "*/")
		if i + 1 == j {
			t := j + 1
			j = strings.Index(s[t:], "*/")
			if j >= 0 {
				j += t
			}
		}
		if i >= 0 && j >= 0 {
			line += s[0:i]
			if len(s) > j + 2 {
				line += s[j + 2:]
			}
			s = line
			line = ""
		} else if i >= 0 {
			line += s[0:i]
			inComment = true
			continue
		} else if inComment && j >= 0 {
			if len(s) > j + 2 {
				line += s[j + 2:]
			}
			s = line
			line = ""
			inComment = false
		}

		if inComment || s == "" {
			continue
		}
		res = append(res, s)
	}

	return res
}

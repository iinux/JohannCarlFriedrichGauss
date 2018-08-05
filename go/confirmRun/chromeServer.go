package main

import (
	"fmt"
	"net/http"
	"log"
	"os/exec"
	"io/ioutil"
)

func runCmdHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprint(w, "Forbidden!")
		return
	}
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	if len(r.Form["key"]) < 1 || r.Form["key"][0] != "911" {
		fmt.Fprint(w, "Forbidden!")
		return
	}

	program := "/opt/google/chrome/chrome"
	cmd := exec.Command(program, r.Form["args"]...)
	runCmd(cmd)

	fmt.Fprintf(w, "args = %s", r.Form["args"])
}

func runCmd(cmd *exec.Cmd) string {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return string(opBytes)
}

func main() {
	http.HandleFunc("/", runCmdHandle)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

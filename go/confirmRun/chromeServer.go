package main

import (
	"fmt"
	"net/http"
	"log"
	"os/exec"
	"io/ioutil"
	"github.com/go-vgo/robotgo/clipboard"
)

func runCmdHandle(w http.ResponseWriter, r *http.Request) {
	if auth(w, r) == false {
		return
	}

	// need system env var
	program := "/opt/google/chrome/chrome"
	// for mac
	// program := "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
	cmd := exec.Command(program, r.Form["args"]...)
	runCmd(cmd)

	fmt.Fprintf(w, "args = %s", r.Form["args"])
}

func auth(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "POST" {
		fmt.Fprint(w, "Forbidden!")
		return false
	}
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	if len(r.Form["key"]) < 1 || r.Form["key"][0] != "911" {
		fmt.Fprint(w, "Forbidden!")
		return false
	}
	return true
}

func getClipHandle(w http.ResponseWriter, r *http.Request) {
	if auth(w, r) == false {
		return
	}

	// need system env var
	/*
	program := "/usr/bin/xclip"
	cmd := exec.Command(program, "-o")
	text := runCmd(cmd)
	*/
	text, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Fprint(w, text)

}

func setClipHandle(w http.ResponseWriter, r *http.Request) {
	if auth(w, r) == false {
		return
	}

	clipText := r.Form["text"][0]
	/*
	copyCmd := exec.Command("xclip", "-in", "-selection", "clipboard")
	in, err := copyCmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := copyCmd.Start(); err != nil {
		fmt.Println(err)
		return
	}
	if _, err := in.Write([]byte(clipText)); err != nil {
		fmt.Println(err)
		return
	}
	if err := in.Close(); err != nil {
		fmt.Println(err)
		return
	}
	copyCmd.Wait()
	*/
	clipboard.WriteAll(clipText)
	fmt.Fprintf(w, "set %s OK", clipText)
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
	http.HandleFunc("/get-clip", getClipHandle)
	http.HandleFunc("/set-clip", setClipHandle)
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

package main

import (
	"log"
	"net/http"
	"os"
	"html/template"
	"encoding/json"
	"encoding/xml"
	"path"
	//for extracting service credentials from VCAP_SERVICES
	//"github.com/cloudfoundry-community/go-cfenv"
)

type Profile struct {
	Name string
	Hobbies []string
}
type ProfileXml struct {
	Name string
	Hobbies []string `xml:"Hobbies>Hobby"`
}

const (
	DEFAULT_PORT = "8080"
	DEFAULT_HOST = "localhost"
)

var index = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))

func helloWorld(w http.ResponseWriter, req *http.Request) {
	index.Execute(w, nil)
}

func jsonFormat(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//w.Write([]byte("OK"))
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}
	js, err := json.Marshal(profile)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(js)
	w.WriteHeader(200)
}
func xmlFormat(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/xml")
	//w.Write([]byte("OK"))
	profile := ProfileXml{"Alex", []string{"snowboarding", "programming"}}
	x, err := xml.Marshal(profile)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(x)
	w.WriteHeader(200)
}
func ServingFile(w http.ResponseWriter, req *http.Request)  {
	fp := path.Join("static", "images", "newapp-icon.png")
	http.ServeFile(w, req, fp)
}


func main() {
	var port string
	if port = os.Getenv("VCAP_APP_PORT"); len(port) == 0 {
		port = DEFAULT_PORT
	}

	var host string
	if host = os.Getenv("VCAP_APP_HOST"); len(host) == 0 {
		host = DEFAULT_HOST
	}

	http.HandleFunc("/", helloWorld)
	http.HandleFunc("/1.json", jsonFormat)
	http.HandleFunc("/1.xml", xmlFormat)
	http.HandleFunc("/servingFile", ServingFile)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("Starting app on %+v:%+v\n", host, port)
	http.ListenAndServe(host + ":" + port, nil)
}

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Get(url string, data url.Values) {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}
	fmt.Println(string(body))
}

func Post(requestUrl string, data url.Values) {
	resp, err := http.PostForm(requestUrl, data)

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}

func sayHelloName(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()   //解析参数，默认是不会解析的
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	_, _ = fmt.Fprintf(w, "Hello %s !", r.Form["name"]) //这个写入到w的是输出到客户端的
}

func main() {
	http.HandleFunc("/", sayHelloName)       //设置访问的路由
	go func() {
		err := http.ListenAndServe(":8887", nil) //设置监听的端口
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()
	requestMethod := "Post"
	requestUrl := "http://www.01happy.com/demo/accept.php"
	requestData := url.Values{"key": {"Value"}, "id": {"123"}}
	functionMap := map[string]func(string, url.Values) {
		"Get": Get,
		"Post": Post,
	}
	doFunction := functionMap[requestMethod]
	doFunction(requestUrl, requestData)
}

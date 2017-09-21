package main

import (
    "fmt"
    "net/http"
    "log"
    "os/exec"
    "io/ioutil"
)

func gitPull(w http.ResponseWriter, r *http.Request) {
    // 执行系统命令
    // 第一个参数是命令名称
    // 后面参数可以有多个，命令参数
    cmd := exec.Command("/usr/bin/git", "pull")
    // 获取输出对象，可以从该对象中读取输出结果
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    // 保证关闭输出流
    defer stdout.Close()
    // 运行命令
    if err := cmd.Start(); err != nil {
        log.Fatal(err)
    }
    // 读取输出结果
    opBytes, err := ioutil.ReadAll(stdout)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Fprint(w, string(opBytes))
    fmt.Fprint(w, "OK")
}

func main() {
    http.HandleFunc("/", gitPull)       //设置访问的路由
    err := http.ListenAndServe(":9920", nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}


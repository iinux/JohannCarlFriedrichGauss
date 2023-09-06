package main

import (
    "fmt"
    "context"
    "log"

    "github.com/tsuna/gohbase"
    "github.com/tsuna/gohbase/hrpc"
)

func main() {
    // 创建 HBase 连接
    client := gohbase.NewClient("127.0.0.1:2181") // 根据你的 HBase 部署设置主机和端口

    myTable := "nn:t2"

    // 构建 Put 请求
    putRequest, err := hrpc.NewPutStr(context.TODO(), myTable, "row-key", map[string]map[string][]byte{
        "f1": {
            "qualifier1": []byte("value1"),
            "qualifier2": []byte("value2"),
        },
    })

    if err != nil {
        log.Fatalf("Put request creation failed: %v", err)
    }

    // 发送 Put 请求
    _, err = client.Put(putRequest)
    if err != nil {
        log.Fatalf("Put request failed: %v", err)
    } else {
        fmt.Println("Put request successful")
    }

    // 构建 Get 请求
    getRequest, err := hrpc.NewGetStr(context.TODO(), myTable, "row-key")

    if err != nil {
        log.Fatalf("Get request creation failed: %v", err)
    }

    // 发送 Get 请求
    getResponse, err := client.Get(getRequest)
    if err != nil {
        log.Fatalf("Get request failed: %v", err)
    } else {
        for _, cell := range getResponse.Cells {
            fmt.Printf("Family: %s, Qualifier: %s, Value: %s\n",
                cell.Family, cell.Qualifier, cell.Value)
        }
    }

    // 关闭 HBase 连接
    client.Close()
}


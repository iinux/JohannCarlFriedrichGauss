package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"flag"
	"net"
)

type Header struct {
	Length int32
	Type   int32
	Uid    int32
	Serid  int32
}

type Request struct {
	Service string
	Call    string
	Params  []string
	Env     map[string]string
}

type Response struct {
	Errno int32  `json:"errno"`
	Msg   string `json:"msg"`
	Data  string `json:"data"`
}

var (
	headSize  = 16
	pkgMaxLen = 10240 // 最大包体长度
	vrfUser   = true
)

var (
	addr = flag.String("addr", "127.0.0.1:9999", "TCP address to listen to")
)

func main() {
	flag.Parse()

	var tcpAddr *net.TCPAddr
	tcpAddr, err := net.ResolveTCPAddr("tcp", *addr)
	if err != nil {
        fmt.Println("Error to resolve addr: ", err)
		return
    }
	tcpListener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
        fmt.Println("Error to linsten tcp: ", err)
		return
    }

	defer tcpListener.Close()

	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			continue
		}

		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn)
	}

}

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("disconnected: " + ipStr)
		conn.Close()
	}()

	// 读取数据
	buf := make([]byte, pkgMaxLen)
	length, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error to read message because of: ", err)
		return
	}
	if length == pkgMaxLen || length < headSize {
		fmt.Println("Got error packet length: ", length)
		return
	}

	// 解析包头
	var head Header
	pack := bytes.NewBuffer(buf[:headSize])
	if err = binary.Read(pack, binary.BigEndian, &head); err != nil {
		fmt.Println("Error to read message because of: ", err)
		return
	}

	// 解析包体
	var req Request
	if err = json.Unmarshal(buf[headSize:length], &req); err != nil {
		fmt.Println("Error to parse request because of: ", err)
		return
	}

	// 响应信息
	res := call(&req, &head)
	rspBody, err := json.Marshal(res)
	if err != nil {
		fmt.Println("json encode err: ", err)
	}

	var headBuf bytes.Buffer
	rspHead := Header{
		Length: int32(len(rspBody)),
		Type:   head.Type,
		Uid:    head.Uid,
		Serid:  head.Serid,
	}
	binary.Write(&headBuf, binary.BigEndian, rspHead)

	rsp := append(headBuf.Bytes(), rspBody...)
	conn.Write(rsp)

	// reader := bufio.NewReader(conn)

	// for {
	//     message, err := reader.ReadString('\n')
	//     if err != nil {
	// 		fmt.Println("an error occourd")
	//         return
	//     }

	//     fmt.Println(string(message))
	//     msg := time.Now().String() + "\n"
	//     b := []byte(msg)
	//     conn.Write(b)
	// }
}

func verifyUser(user, pwd string) bool {
	if user == "echeng" && pwd == "123456" {
		return true
	}
	return false
}

func call(req *Request, head *Header) Response {
	// 请求参数是否正确
	if len(req.Service) == 0 || len(req.Call) == 0 {
		return Response{
			Errno: 10010,
			Msg:   "invalid service or call",
			Data:  "",
		}
	}

	// 侦测服务器是否存活
	if req.Call == "PING" {
		return Response{
			Errno: 10010,
			Msg:   "PONG",
			Data:  "",
		}
	}

	// 验证密码是否正确
	if vrfUser {
		env := req.Env
		if len(env["user"]) == 0 || len(env["password"]) == 0 {
			return Response{
				Errno: 10010,
				Msg:   "",
				Data:  "unauthorized",
			}
		}

		if !verifyUser(env["user"], env["password"]) {
			return Response{
				Errno: 10010,
				Msg:   "",
				Data:  "unauthorized",
			}
		}
	}

	// 调用具体业务逻辑
	// .............
	return Response{
		Errno: 0,
		Msg:   "",
		Data:  `{"lang":"ch", "content":"1233456"}`,
	}
}

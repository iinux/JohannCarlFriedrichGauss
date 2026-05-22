// af_alg_aes.go
//
// 使用 Linux 内核加密接口 AF_ALG 进行 AES-CBC 加解密。
//
// 要求:
//   - Linux 内核 >= 2.6.38 (启用 CONFIG_CRYPTO_USER_API_SKCIPHER)
//   - golang.org/x/sys/unix
//
// 运行:
//   go mod init af_alg_demo
//   go get golang.org/x/sys/unix
//   go run af_alg_aes.go

package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// rawAccept 直接调用 accept(2)。某些内核版本对 AF_ALG 上的 accept4 + SOCK_CLOEXEC
// 处理有缺陷,会返回 ECONNABORTED,因此这里走原始 accept。
func rawAccept(fd int) (int, error) {
	r, _, errno := syscall.Syscall(syscall.SYS_ACCEPT, uintptr(fd), 0, 0)
	if errno != 0 {
		return -1, errno
	}
	return int(r), nil
}

const blockSize = 16

// pkcs7Pad 按 PKCS#7 规则填充到 blockSize 倍数。
func pkcs7Pad(data []byte, bs int) []byte {
	padLen := bs - len(data)%bs
	return append(data, bytes.Repeat([]byte{byte(padLen)}, padLen)...)
}

// pkcs7Unpad 去掉 PKCS#7 填充。
func pkcs7Unpad(data []byte) ([]byte, error) {
	n := len(data)
	if n == 0 {
		return nil, fmt.Errorf("空数据无法去填充")
	}
	padLen := int(data[n-1])
	if padLen < 1 || padLen > blockSize || padLen > n {
		return nil, fmt.Errorf("PKCS7 填充无效")
	}
	for i := n - padLen; i < n; i++ {
		if int(data[i]) != padLen {
			return nil, fmt.Errorf("PKCS7 填充无效")
		}
	}
	return data[:n-padLen], nil
}

// afAlgCrypt 通过 AF_ALG 执行对称加/解密。
//
//	op:      unix.ALG_OP_ENCRYPT 或 unix.ALG_OP_DECRYPT
//	algName: 内核算法名,例如 "cbc(aes)"
//	data:    输入长度必须是 blockSize 整数倍
func afAlgCrypt(op uint32, algName string, key, iv, data []byte) ([]byte, error) {
	// 1. 创建算法 socket
	algFd, err := unix.Socket(unix.AF_ALG, unix.SOCK_SEQPACKET, 0)
	if err != nil {
		return nil, fmt.Errorf("socket(AF_ALG): %w", err)
	}
	defer unix.Close(algFd)

	// 2. bind 到具体算法
	if err := unix.Bind(algFd, &unix.SockaddrALG{Type: "skcipher", Name: algName}); err != nil {
		return nil, fmt.Errorf("bind(skcipher,%s): %w", algName, err)
	}

	// 3. 设置密钥
	if err := unix.SetsockoptString(algFd, unix.SOL_ALG, unix.ALG_SET_KEY, string(key)); err != nil {
		return nil, fmt.Errorf("setsockopt ALG_SET_KEY: %w", err)
	}

	// 4. accept 得到操作 fd
	opFd, err := rawAccept(algFd)
	if err != nil {
		return nil, fmt.Errorf("accept: %w", err)
	}
	defer unix.Close(opFd)

	// 5. 构造控制消息(cmsg): ALG_SET_OP + ALG_SET_IV
	opPayload := make([]byte, 4)
	binary.NativeEndian.PutUint32(opPayload, op)

	ivPayload := make([]byte, 4+len(iv)) // struct af_alg_iv { __u32 ivlen; __u8 iv[0]; }
	binary.NativeEndian.PutUint32(ivPayload[:4], uint32(len(iv)))
	copy(ivPayload[4:], iv)

	oob := append(buildCmsg(unix.SOL_ALG, unix.ALG_SET_OP, opPayload),
		buildCmsg(unix.SOL_ALG, unix.ALG_SET_IV, ivPayload)...)

	// 6. 一次性发送 op + iv (cmsg) 和数据
	if err := unix.Sendmsg(opFd, data, oob, nil, 0); err != nil {
		return nil, fmt.Errorf("sendmsg: %w", err)
	}

	// 7. 读取结果(长度与输入相同)
	out := make([]byte, len(data))
	n, err := unix.Read(opFd, out)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}
	if n != len(data) {
		return nil, fmt.Errorf("短读: %d != %d", n, len(data))
	}
	return out, nil
}

// buildCmsg 构造一条 cmsghdr + payload(按 cmsg 对齐填充)。
func buildCmsg(level, typ int32, payload []byte) []byte {
	buf := make([]byte, unix.CmsgSpace(len(payload)))
	h := (*unix.Cmsghdr)(unsafe.Pointer(&buf[0]))
	h.Level = level
	h.Type = typ
	h.SetLen(unix.CmsgLen(len(payload)))
	copy(buf[unix.CmsgLen(0):], payload)
	return buf
}

func encrypt(key, iv, plaintext []byte) ([]byte, error) {
	return afAlgCrypt(unix.ALG_OP_ENCRYPT, "cbc(aes)", key, iv, pkcs7Pad(plaintext, blockSize))
}

func decrypt(key, iv, ciphertext []byte) ([]byte, error) {
	if len(ciphertext)%blockSize != 0 {
		return nil, fmt.Errorf("密文长度必须是 %d 的倍数", blockSize)
	}
	padded, err := afAlgCrypt(unix.ALG_OP_DECRYPT, "cbc(aes)", key, iv, ciphertext)
	if err != nil {
		return nil, err
	}
	return pkcs7Unpad(padded)
}

func main() {
	key := make([]byte, 32) // AES-256
	iv := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		log.Fatal(err)
	}
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}

	plaintext := []byte("AF_ALG 内核加密接口测试: hello, 世界!")
	fmt.Printf("明文 (%d 字节): %q\n", len(plaintext), plaintext)
	fmt.Println("密钥 (hex):", hex.EncodeToString(key))
	fmt.Println("IV   (hex):", hex.EncodeToString(iv))

	ct, err := encrypt(key, iv, plaintext)
	if err != nil {
		log.Fatalf("加密失败: %v", err)
	}
	fmt.Println("密文 (hex):", hex.EncodeToString(ct))

	pt, err := decrypt(key, iv, ct)
	if err != nil {
		log.Fatalf("解密失败: %v", err)
	}
	fmt.Printf("解密 (%d 字节): %q\n", len(pt), pt)

	if !bytes.Equal(pt, plaintext) {
		log.Fatal("解密结果与原文不一致!")
	}
	fmt.Println("OK: 加解密往返一致")
}

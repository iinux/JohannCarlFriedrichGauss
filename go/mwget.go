package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	BufferSize = 1024
	ThreadNum = 4  // 默认使用4个线程进行下载
)

// Chunk 表示要下载的文件块信息
type Chunk struct {
	Start int64
	End   int64
	ID    int
	Path  string
	URL   string
}

// MultiThreadDownloader 多线程下载器结构体
type MultiThreadDownloader struct {
	URL      string
	Filename string
	Threads  int
	client   *http.Client
}

// NewMultiThreadDownloader 创建一个新的多线程下载器实例
func NewMultiThreadDownloader(url, filename string, threads int) *MultiThreadDownloader {
	if threads <= 0 {
		threads = ThreadNum
	}
	return &MultiThreadDownloader{
		URL:      url,
		Filename: filename,
		Threads:  threads,
		client:   &http.Client{},
	}
}

// GetFileSize 获取远程文件大小
func (mtd *MultiThreadDownloader) GetFileSize() (int64, error) {
	resp, err := mtd.client.Head(mtd.URL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("无法获取文件大小")
	}

	size, err := strconv.ParseInt(contentLength, 10, 64)
	if err != nil {
		return 0, err
	}

	return size, nil
}

// Download 启动多线程下载
func (mtd *MultiThreadDownloader) Download() error {
	fileSize, err := mtd.GetFileSize()
	if err != nil {
		return fmt.Errorf("获取文件大小失败: %v", err)
	}

	// 计算每个线程需要下载的数据量
	chunkSize := fileSize / int64(mtd.Threads)

	var chunks []Chunk
	for i := 0; i < mtd.Threads; i++ {
		start := int64(i) * chunkSize
		end := start + chunkSize - 1

		// 最后一个分片需要包含剩余的所有数据
		if i == mtd.Threads-1 {
			end = fileSize - 1
		}

		chunk := Chunk{
			Start: start,
			End:   end,
			ID:    i,
			Path:  fmt.Sprintf("%s.part%d", mtd.Filename, i),
			URL:   mtd.URL,
		}
		chunks = append(chunks, chunk)
	}

	// 创建临时文件用于合并
	tempFile, err := os.Create(mtd.Filename + ".tmp")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// 扩展文件到目标大小
	if err := tempFile.Truncate(fileSize); err != nil {
		return err
	}

	// 使用WaitGroup等待所有goroutine完成
	var wg sync.WaitGroup

	// 启动多个goroutine下载各个分片
	for _, chunk := range chunks {
		wg.Add(1)
		go func(c Chunk) {
			defer wg.Done()
			err := mtd.downloadChunk(tempFile, c)
			if err != nil {
				fmt.Printf("下载分片 %d 时出错: %v\n", c.ID, err)
			} else {
				fmt.Printf("分片 %d 下载完成 [%d-%d]\n", c.ID, c.Start, c.End)
			}
		}(chunk)
	}

	// 等待所有goroutine完成
	wg.Wait()

	// 检查是否所有部分都已下载完成
	finalSize, err := mtd.GetFileSize()
	if err != nil {
		return err
	}

	tempStat, err := tempFile.Stat()
	if err != nil {
		return err
	}

	if tempStat.Size() != finalSize {
		return fmt.Errorf("下载不完整，期望大小: %d, 实际大小: %d", finalSize, tempStat.Size())
	}

	// 将临时文件重命名为最终文件名
	tempFile.Close()
	err = os.Rename(mtd.Filename+".tmp", mtd.Filename)
	if err != nil {
		return err
	}

	fmt.Println("下载完成:", mtd.Filename)
	return nil
}

// downloadChunk 下载单个分片并写入临时文件
func (mtd *MultiThreadDownloader) downloadChunk(file *os.File, chunk Chunk) error {
	req, err := http.NewRequest("GET", chunk.URL, nil)
	if err != nil {
		return err
	}

	// 设置Range请求头，指定下载范围
	rangeHeader := fmt.Sprintf("bytes=%d-%d", chunk.Start, chunk.End)
	req.Header.Set("Range", rangeHeader)

	resp, err := mtd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 在指定位置开始写入文件
	_, err = file.WriteAt([]byte{}, chunk.Start) // 预分配空间
	if err != nil {
		return err
	}

	// 将数据写入指定的位置
	_, err = io.Copy(&SectionWriter{file, chunk.Start}, resp.Body)
	return err
}

// SectionWriter 包装os.File，实现从指定偏移量开始写入
type SectionWriter struct {
	file   *os.File
	offset int64
}

func (sw *SectionWriter) Write(p []byte) (n int, err error) {
	defer func() {
		sw.offset += int64(n)
	}()

	n, err = sw.file.WriteAt(p, sw.offset)
	return
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: go run multithreaded-wget.go <URL> [filename] [threads]")
		return
	}

	url := os.Args[1]
	filename := filepath.Base(url) // 默认使用URL的最后一段作为文件名
	threads := ThreadNum

	if len(os.Args) > 2 {
		filename = os.Args[2]
	}
	if len(os.Args) > 3 {
		if t, err := strconv.Atoi(os.Args[3]); err == nil && t > 0 {
			threads = t
		}
	}

	downloader := NewMultiThreadDownloader(url, filename, threads)

	fmt.Printf("开始下载: %s\n", url)
	fmt.Printf("保存为: %s\n", filename)
	fmt.Printf("使用线程数: %d\n", threads)

	err := downloader.Download()
	if err != nil {
		fmt.Printf("下载失败: %v\n", err)
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type writeCounter struct {
	Total uint64
	size  int
}

func (f *fileData) InitIPData(url string, path string, size int) (rs interface{}) {
	var tmpData []byte

	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 IP 库")
		err = downloadFile(path, url, size)
		if err != nil {
			rs = err
			return
		}
		log.Printf("已将最新 IP 库保存到本地 %s ", path)
	}

	f.Path, err = os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		rs = err
		return
	}
	defer f.Path.Close()

	tmpData, err = ioutil.ReadAll(f.Path)
	if err != nil {
		log.Println(err)
		rs = err
		return
	}

	f.Data = tmpData
	return true
}

// https://gist.github.com/albulescu/e61979cc852e4ee8f49c
func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.Total += uint64(n)
	wc.printProgress()
	return n, nil
}

func (wc writeCounter) printProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %d KB / %d KB", wc.Total/1024, wc.size)
}

func downloadFile(filepath string, url string, size int) error {
	out, err := os.Create(filepath + ".tmp")
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	counter := &writeCounter{size: size}
	_, err = io.Copy(out, io.TeeReader(resp.Body, counter))
	if err != nil {
		return err
	}

	fmt.Print("\n")

	err = os.Rename(filepath+".tmp", filepath)
	if err != nil {
		return err
	}

	return nil
}

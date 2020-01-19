package main

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
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

// github.com/freshcn/qqwry/blob/bb0d1d8ade3948d506ae836c641c2cbe0ad2ca45/download.go
func getKey() (uint32, error) {
	resp, err := http.Get("https://qqwry.mirror.noc.one/copywrite.rar")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(body[5*4:]), nil
}

func (f *fileData) InitIPData(url string, path string, size int) (rs interface{}) {
	var tmpData []byte

	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Println("文件不存在，尝试从网络获取最新 IP 库")
		if path != "ipv4.dat" {
			err = downloadFile(path, url, size)
			if err != nil {
				rs = err
				return
			}
		} else {
			err = downloadFile("encrypted.tmp", url, size)
			if err != nil {
				rs = err
				return
			}
			if err = decrypt(); err != nil {
				rs = err
				return
			}
		}
		log.Printf("已将最新 IP 库保存到本地")
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

func decrypt() error {
	File, err := os.OpenFile("encrypted.tmp", os.O_RDONLY, 0400)
	if err != nil {
		return err
	}
	defer File.Close()

	if body, err := ioutil.ReadAll(File); err == nil {
		if key, err := getKey(); err == nil {
			for i := 0; i < 0x200; i++ {
				key = key * 0x805
				key++
				key = key & 0xff
				body[i] = byte(uint32(body[i]) ^ key)
			}

			reader, err := zlib.NewReader(bytes.NewReader(body))
			if err != nil {
				return err
			}

			Data, err := ioutil.ReadAll(reader)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile("ipv4.dat", Data, 0644); err == nil {
				_ = os.Remove("encrypted.tmp")
			}
		}
	}
	return err
}

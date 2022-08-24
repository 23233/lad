package lad

import (
	"bufio"
	"bytes"
	"context"
	"embed"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	reqClient = http.Client{
		Timeout: 20 * time.Second,
	}
	GetRemoteFail = errors.New("获取远程文件失败")
)

// Load 路径加载本地文件
func (ac *acMachine) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		ac.Add(line)
	}
}

// LoadOfFolder 加载文件夹下所有文件
func (ac *acMachine) LoadOfFolder(folder embed.FS) error {
	err := fs.WalkDir(folder, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		if d.IsDir() {
			return nil
		}
		// 仅支持.txt .data .dict后缀名的文件
		allowSuffix := []string{".txt", ".data", ".dict"}
		for _, suffix := range allowSuffix {
			if strings.HasSuffix(d.Name(), suffix) {
				b, err := fs.ReadFile(folder, path)
				if err != nil {
					return err // or panic or ignore
				}
				return ac.loadOfByte(b)
			}
		}
		return nil
	})
	return err
}

func (ac *acMachine) loadOfByte(file []byte) error {
	reader := bufio.NewReader(bytes.NewReader(file))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		ac.Add(line)
	}
}

// LoadRemote 加载远程文件 未验证内容格式是否正确 请自行判断
func (ac *acMachine) LoadRemote(url string, timeout time.Duration) error {
	ctx, _ := context.WithTimeout(context.TODO(), timeout)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return GetRemoteFail
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return ac.loadOfByte(body)
}

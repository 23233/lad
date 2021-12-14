package lad

import (
	"embed"
	"testing"
	"time"
)

//go:embed data/*
var folder embed.FS

func TestLoadOfFolder(t *testing.T) {
	machine := New()
	err := machine.LoadOfFolder(folder)
	if err != nil {
		t.Fatal(err)
	}
	machine.root.view()
}

func TestLoadRemote(t *testing.T) {
	machine := New()
	// 若存在墙的问题 按如下解决
	// 1. https://ipaddress.com/website/raw.githubusercontent.com 获取到ip
	// 2. 设置hosts 并生效 window修改hosts文件会直接生效 linux需要刷新dns缓存
	err := machine.LoadRemote("https://raw.githubusercontent.com/23233/lad/master/test.data", 20*time.Second)
	if err != nil {
		t.Fatal(err)
	}

	machine.root.view()
}

func TestLoader(t *testing.T) {
	t.Run("测试本地文件加载", TestLoadOfFolder)
	t.Run("测试远程文件加载", TestLoadRemote)
}

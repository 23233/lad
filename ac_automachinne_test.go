package lad

import (
	"testing"
)

func TestAcMachine_Add(t *testing.T) {
	machine := New()
	machine.Add("我们")
	machine.Add("历史")
	machine.Add("listen to me")
	machine.AddOfList([]string{"哎哟", "走起"})
	p := machine.root

	p.view()
}

func TestAcMachine_Load(t *testing.T) {
	machine := New()
	if err := machine.Load("./test.data"); err != nil {
		t.Error(err)
	}

	machine.root.view()
}

func TestAcMachine_Build(t *testing.T) {
	machine := New()
	if err := machine.Load("./test.data"); err != nil {
		t.Error(err)
	}
	machine.Build()
	machine.root.view()
}

func TestAcMachine_Find(t *testing.T) {
	machine := New()
	if err := machine.Load("./test.data"); err != nil {
		t.Error(err)
	}
	machine.Build()

	t.Log(machine.Find("ab    阿宾   cdadfadfadfadf")) // [阿宾]
}

func TestAcMachine_Match(t *testing.T) {
	machine := New()
	if err := machine.Load("./test.data"); err != nil {
		t.Error(err)
	}
	machine.Build()

	t.Log(machine.Match("abx     阿宾  cdadfadfadfadf")) // true
}

func TestAcMachine_Replace(t *testing.T) {
	machine := New()
	if err := machine.Load("./test.data"); err != nil {
		t.Error(err)
	}
	machine.Build()

	t.Log(machine.Replace("ab    阿宾   cdadfadfadfadf", "****")) // ab    ****   cdadfadfadfadf
}

func TestLab(t *testing.T) {
	t.Run("新增", TestAcMachine_Add)
	t.Run("路径文件加载", TestAcMachine_Load)
	t.Run("状态机构建", TestAcMachine_Build)
	t.Run("查找匹配", TestAcMachine_Find)
	t.Run("查找替换", TestAcMachine_Match)
	t.Run("查找替换", TestAcMachine_Replace)
}

func BenchmarkAcMachine_Find(b *testing.B) {
	machine := New()
	if err := machine.Load("./test.data"); err != nil {
		b.Error(err)
	}
	machine.Build()
	b.Run("find", func(b *testing.B) {

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				machine.Find("cdadfadfadfadfab    阿宾   cdadfadfadfadf，俣哈萨克斯坦喝茶；罪域国；工；国；甘；adfa;dflklazdlfjaldf a:jfla工；期刊")
			}
		})
	})
}

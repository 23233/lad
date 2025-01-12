package lad

import (
	"container/list"
	"fmt"
	"strings"
)

const rootRaw = "/"

type acNode struct {
	raw      string
	children map[string]*acNode
	isEnd    bool
	length   int
	fail     *acNode
}

func NewAcNode(raw string) *acNode {
	return &acNode{
		raw:      raw,
		children: make(map[string]*acNode),
	}
}

func (an *acNode) view() {
	fmt.Println(an.raw, ",", an.length, ",fail:", an.fail)
	if len(an.children) == 0 {
		return
	}
	for _, node := range an.children {
		node.view()
	}
}

type AcMachine struct {
	root *acNode
}

func New() *AcMachine {
	root := NewAcNode(rootRaw)
	return &AcMachine{
		root: root,
	}
}

// add 添加模式串
func (ac *AcMachine) add(pattern string) {
	p := ac.root
	tok := newToken(pattern)
	length := 0
	for str := tok.next(); str != ""; str = tok.next() {
		if _, ok := p.children[str]; !ok {
			newNode := NewAcNode(str)
			p.children[str] = newNode
		}
		p = p.children[str]
		length++
	}

	p.length = length
	p.isEnd = true
}

// Add 新增单个词
func (ac *AcMachine) Add(word string) {
	ac.add(word)
}

// AddOfList 新增一堆词
func (ac *AcMachine) AddOfList(word []string) {
	for i := 0; i < len(word); i++ {
		ac.add(word[i])
	}
}

// Build 构建自动机
func (ac *AcMachine) Build() {
	l := list.New()

	l.PushBack(ac.root)

	for l.Len() > 0 {
		e := l.Front()
		l.Remove(e)
		p := e.Value.(*acNode)

		for _, pc := range p.children {
			if p == ac.root {
				pc.fail = ac.root
			} else {
				q := p.fail

				for q != nil {
					if qc, ok := q.children[pc.raw]; ok {
						pc.fail = qc
						break
					}
					q = q.fail
				}

				if q == nil {
					pc.fail = ac.root
				}
			}
			l.PushBack(pc)
		}
	}
}

// Find 查找
func (ac *AcMachine) Find(text string) []string {
	rs := make([]string, 0)
	ac.match(text, func(tok *token, node *acNode) {
		rs = append(rs, tok.prevNStr(tok.index, node.length))
	})

	return rs
}

// Match 匹配
func (ac *AcMachine) Match(text string) bool {
	rs := false
	ac.match(text, func(tok *token, node *acNode) {
		rs = true
	})
	return rs
}
func (ac *AcMachine) match(text string, fn func(tok *token, node *acNode)) {
	p := ac.root
	tok := newToken(text)
	for {
		str := tok.next()
		if str == "" {
			break
		}
		for {
			if _, ok := p.children[str]; !ok && p != ac.root {
				p = p.fail
				continue
			}
			break
		}

		p = p.children[str]

		if p == nil {
			p = ac.root
		}

		tmp := p
		for tmp != ac.root {
			if tmp.isEnd {
				fn(tok, tmp)
			}
			tmp = tmp.fail
		}
	}
}

// Replace 替换
func (ac *AcMachine) Replace(text, target string) string {
	rs := ""
	ac.match(text, func(tok *token, node *acNode) {
		if rs == "" {
			rs = string(tok.origin)
		}
		rs = strings.Replace(rs, tok.prevNStr(tok.index, node.length), target, -1)
	})
	return rs
}

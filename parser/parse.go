package parser

import (
	"log"
	"strings"
)

type Source struct {
	Description
}

type Description struct {
	ModuleDef interface{}
}

type Paramlist struct{}

type Portlist struct {
	Port
}

type Port struct {
	Name string
	// Noneになっているのがわからない
}

type Wire struct {
	Name string
	// Falseになっているのがわからない
}

type Input struct {
	Name string
	// Falseになっているのがわからない
}

type Output struct {
	Name string
	// Falseになっているのがわからない
}

type Assign struct {
	Lvalue interface{}
	Rvalue interface{}
}

type Identifier struct {
	Name string
}

type GateType int

const (
	_ GateType = iota
	Xor
	And
	Or
)

// テスト
type At int

var TestMap map[At][]Node
var PortList []string

type Node struct {
	In1  string
	In2  string
	Gate GateType
	Out  string
}

func NewParse() {
	file := NewReader("test.txt")
	scanner := NewScanner(file)
	for scanner.Scan() {
		txt := scanner.Text()
		strline := strings.Fields(txt)

		switch strline[0] {
		case "Assign:":
			assignwidth := strings.Fields(strline[1])
			from := assignwidth[1]
			to := assignwidth[3]
			if from != to[:len(to)-1] {
				// 一次的にこうする　もしかしたらAssignでto, fromが複数もあるかもしれない
				log.Fatalf("from and to not equal Error: from:%s, to:%s", from, to)
			}
		case "Lvalue:":

		case "Rvalue:":
		case "Identifier:":

		case "Xor:":

		case "And:":
		case "Or:":

		default:
			continue
		}
	}
}

// (from N to N)をパースしてNを取得
func parseFromTo() {

}

// (at N)をパースしてNを取得
func parseAt() {

}

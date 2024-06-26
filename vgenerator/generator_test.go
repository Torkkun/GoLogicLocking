package vgenerator

import (
	"fmt"
	"goll/parser/yosysjson"
	"goll/vgenerator/funcmap"
	"log"
	"os"
	"testing"
	"text/template"
)

func TestStringTest(t *testing.T) {

	text := "module {{ .ModuleName }} ({{ argElement .PortList }});"
	tpl := template.New("")
	mname := "test_module"
	list := []string{"A", "B", "C", "D"}

	val := struct {
		ModuleName string
		PortList   []string
	}{
		ModuleName: mname,
		PortList:   list,
	}

	funcmap := template.FuncMap{
		"argElement": func(io []string) string {
			var arg string
			for i, v := range io {
				if len(io)-1 == i {
					arg += v
				} else {
					arg += fmt.Sprintf("%s, ", v)
				}
			}
			return arg
		},
	}

	tpl = tpl.Funcs(funcmap)
	var err error

	tpl, err = tpl.Parse(text)
	if err != nil {
		log.Fatal(err)
	}

	if err := tpl.Execute(os.Stdout, val); err != nil {
		log.Fatalln(err)
	}
}

var val = struct {
	ModuleName string
	PortList   []string
	Parameter  []string // not implement
	PortDecl   []*funcmap.PortDecl
	Register   []string              // not implement
	EventDecl  []string              // not implement
	NetDecl    []*funcmap.NetDecl    // not implement
	AssignDecl []*funcmap.AssignDecl // not implement
}{
	ModuleName: "test_module",
	PortList:   []string{"A", "B", "C", "D"},
}

func TestText(t *testing.T) {
	if err := NewGenerator(val, true); err != nil {
		log.Fatalln(err)
	}
}

func TestElement(t *testing.T) {
	val.PortDecl = []*funcmap.PortDecl{
		{
			PortType:   funcmap.Input,
			SignalName: "A",
		},
		{
			PortType:   funcmap.Input,
			SignalName: "B",
		},
		{
			PortType:   funcmap.Input,
			SignalName: "C",
		},
		{
			PortType:   funcmap.Output,
			SignalName: "D",
		},
		{
			PortType: funcmap.Input,
			BitWidth: &funcmap.BitWidth{
				MSB: 4,
				LSB: 0,
			},
			SignalName: "E",
		},
	}
	val.NetDecl = []*funcmap.NetDecl{
		{
			NetType: funcmap.Wire,
			NetName: "_1_",
		},
		{
			NetType: funcmap.Wire,
			BitWidth: &funcmap.BitWidth{
				MSB: 2,
				LSB: 0,
			},
			NetName: "_2_",
		},
	}
	testconn1 := funcmap.Logic{
		Y: "out1",
		A: "in1",
		B: "in2",
	}
	testconn2 := funcmap.BuforNot{
		Y: "out1",
		A: "in1",
	}
	val.AssignDecl = []*funcmap.AssignDecl{
		{
			ExpressionType: yosysjson.AND,
			Connection:     testconn1,
		},
		{
			ExpressionType: yosysjson.NAND,
			Connection:     testconn1,
		},
		{
			ExpressionType: yosysjson.OR,
			Connection:     testconn1,
		},
		{
			ExpressionType: yosysjson.NOR,
			Connection:     testconn1,
		},
		{
			ExpressionType: yosysjson.XOR,
			Connection:     testconn1,
		},
		{
			ExpressionType: yosysjson.XNOR,
			Connection:     testconn1,
		},
		{
			ExpressionType: yosysjson.BUF,
			Connection:     testconn2,
		},
		{
			ExpressionType: yosysjson.NOT,
			Connection:     testconn2,
		},
	}
	if err := NewGenerator(val, false); err != nil {
		log.Fatalln(err)
	}
}

package vgenerator

import (
	"fmt"
	"goll/vgenerator/funcmap"
	"log"
	"os"
	"testing"
	"text/template"
)

func TestStringTest(t *testing.T) {

	text := "module {{ .Module }} ({{ argElement .IO }});"
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
	ModuleName    string
	PortList      []string
	Parameter     []string // not implement
	PortDecl      []*funcmap.PortDecl
	Register      []string // not implement
	EventDecl     []string // not implement
	NetDecl       []string // not implement
	PremitiveDecl []string // not implement
}{
	ModuleName: "test_module",
	PortList:   []string{"A", "B", "C", "D"},
}

func TestText(t *testing.T) {
	if err := NewGenerator(val, true); err != nil {
		log.Fatalln(err)
	}
}

func TestPortListElement(t *testing.T) {
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
	}
	if err := NewGenerator(val, false); err != nil {
		log.Fatalln(err)
	}
}

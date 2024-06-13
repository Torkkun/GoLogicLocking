package vgenerator

import (
	"fmt"
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
		Module string
		IO     []string
	}{
		Module: mname,
		IO:     list,
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

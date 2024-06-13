package funcmap

import (
	"fmt"
	"html/template"
)

func NewFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"argElement": argElement,
	}
	return funcMap
}

func argElement(io []string) string {
	var arg string
	for i, v := range io {
		if len(io)-1 == i {
			arg += v
		} else {
			arg += fmt.Sprintf("%s, ", v)
		}
	}
	return arg
}

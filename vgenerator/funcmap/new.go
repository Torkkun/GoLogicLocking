package funcmap

import "text/template"

func NewFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"portListElement":      portListElement,
		"paramDeclElement":     paramDeclElement,
		"portDeclElement":      portDeclElement,
		"registerDeclElement":  registerDeclElement,
		"eventDeclElement":     eventDeclElement,
		"netDeclElement":       netDeclElement,
		"premitiveDeclElement": premitiveDeclElement,
	}
	return funcMap
}

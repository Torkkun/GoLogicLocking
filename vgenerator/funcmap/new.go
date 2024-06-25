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

type BitWidth struct {
	MSB int
	LSB int
}

type PortType string

const (
	Input  PortType = "input"
	Output PortType = "output"
	InOut  PortType = "inout"
)

type PortDecl struct {
	PortType   PortType
	BitWidth   *BitWidth
	SignalName string
}

type NetType string

const (
	Wire NetType = "wire"
)

type NetDecl struct {
	NetType  NetType // 今のところwireのみ
	BitWidth *BitWidth
	NetName  string
}

package funcmap

import "text/template"

func NewFuncMap() template.FuncMap {
	funcMap := template.FuncMap{
		"portListElement":     portListElement,
		"paramDeclElement":    paramDeclElement,
		"portDeclElement":     portDeclElement,
		"registerDeclElement": registerDeclElement,
		"eventDeclElement":    eventDeclElement,
		"netDeclElement":      netDeclElement,
		"assignDeclElement":   assignDeclElement,
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

type AssignDecl struct {
	ExpressionType string
	Connection     interface{} // Name: input | output
}

type BuforNot struct {
	Y string
	A string
}

type Logic struct {
	Y string
	A string
	B string
}

type Logic3 struct {
	Y string
	A string
	B string
	C string
}

type Logic4 struct {
	Y string
	A string
	B string
	C string
	D string
}

type Mux struct {
	Y string
	S string
	B string
	A string
}

// mux4,8,16

//FF

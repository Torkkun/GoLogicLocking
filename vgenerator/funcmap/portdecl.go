package funcmap

import (
	"fmt"
	"goll/utils"
)

type PortType string

const (
	Input  PortType = "input"
	Output PortType = "output"
)

type PortDecl struct {
	PortType PortType
	BitWidth *struct {
		MSB int
		LSB int
	}
	SignalName string
}

func portDeclElement(portdecls []*PortDecl) string {
	var returnstr string
	for i, portdecl := range portdecls {
		var portdeclstr string
		switch portdecl.PortType {
		case Input:
			portdeclstr = "input"
		case Output:
			portdeclstr = "output"
		}

		if portdecl.BitWidth != nil {
			portdeclstr += utils.WriteSpace()
			portdeclstr += fmt.Sprintf("[%d:%d]", portdecl.BitWidth.MSB, portdecl.BitWidth.LSB)
		}

		portdeclstr += utils.WriteSpace()
		portdeclstr += portdecl.SignalName

		portdeclstr = utils.WriteIndentTab() + portdeclstr
		if i != len(portdecls)-1 {
			portdeclstr = portdeclstr + utils.WriteNewLine()
		}
		returnstr += portdeclstr
	}

	return returnstr
}

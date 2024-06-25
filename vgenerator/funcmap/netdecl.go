package funcmap

import (
	"fmt"
	"goll/utils"
	"log"
)

func netDeclElement(netdecls []*NetDecl) string {
	var returnstr string
	for i, netdecl := range netdecls {
		var netdeclstr string
		switch netdecl.NetType {
		case Wire:
			netdeclstr = "wire"
		default:
			log.Fatalln("This Net Type not undefined")
		}

		if netdecl.BitWidth != nil {
			netdeclstr += utils.WriteSpace()
			netdeclstr += fmt.Sprintf("[%d:%d]", netdecl.BitWidth.MSB, netdecl.BitWidth.LSB)
		}

		netdeclstr += utils.WriteSpace()
		netdeclstr += netdecl.NetName

		netdeclstr = utils.WriteIndentTab() + netdeclstr
		if i != len(netdecls)-1 {
			netdeclstr = netdeclstr + utils.WriteSemiColon() + utils.WriteNewLine()
		} else {
			netdeclstr = netdeclstr + utils.WriteSemiColon()
		}
		returnstr += netdeclstr
	}

	return returnstr
}

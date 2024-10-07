package yosysjson

import (
	"fmt"
	"goll/graph/veryosys"

	mapset "github.com/deckarep/golang-set/v2"
)

func PortsConvertToStoreGraph(ports map[string]Port) []*veryosys.Port {
	var graphPorts []*veryosys.Port
	for name, port := range ports {

		portnumSet := mapset.NewSet[int](port.Bits...)
		if portnumSet.Cardinality() > 1 {
			panic(fmt.Errorf("現在対応しているのはBitNumが一種類のみ: %v", portnumSet.Cardinality()).Error())
		}

		graphPorts = append(graphPorts, &veryosys.Port{
			Direction: port.Direction,
			BitNum:    port.Bits[0],
			BitWidth:  len(port.Bits),
			Name:      name,
		})
	}
	return graphPorts
}

func NetsConvertToStoreGraph(nets map[string]NetName) []*veryosys.NetName {
	var graphNets []*veryosys.NetName

	for netname, net := range nets {
		if net.HideName == 0 {
			continue
		}
		graphNets = append(graphNets, &veryosys.NetName{
			Bits:    net.Bits,
			Netname: netname,
			Attributes: struct{ Src string }{
				net.Attributes.Src,
			},
		})
	}
	return graphNets
}

func CellsConvertToStoreGraph(cells map[string]Cell) ([]*veryosys.Cell, veryosys.Connections) {
	var graphCells []*veryosys.Cell
	graphConns := make(map[string][]*veryosys.Connection)
	for _, cell := range cells {
		if cell.HideName == 0 {
			continue
		}
		var conns []*veryosys.Connection
		for portName, portType := range cell.PortDirections {
			conn := new(veryosys.Connection)
			bitnum := cell.Connections[portName]

			// 接続Bitが1Bitではないときエラー
			if len(bitnum) > 1 {
				err := fmt.Errorf("接続Bitが1Bitではないときエラー:multibit not implement: %v", bitnum)
				panic(err.Error())
			}
			conn.Type = portType
			conn.BitNum = bitnum[0]
			conn.PortName = portName

			conns = append(conns, conn)
		}
		graphCells = append(graphCells, &veryosys.Cell{
			Type: cell.Type,
			Attributes: struct {
				Src string
			}{
				cell.Attributes.Src,
			},
		})
		graphConns[cell.Attributes.Src] = conns
	}
	return graphCells, graphConns
}

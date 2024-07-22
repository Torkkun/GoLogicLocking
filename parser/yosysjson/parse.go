package yosysjson

import (
	"fmt"
	"goll/graph/veryosys"
)

func PortsConvertToStoreGraph(ports map[string]Port) []*veryosys.Port {
	var graphPorts []*veryosys.Port
	for name, port := range ports {
		graphPorts = append(graphPorts, &veryosys.Port{
			Direction: port.Direction,
			BitNum:    port.Bits,
			BitWidth:  len(port.Bits),
			Name:      name,
		})
	}
	return graphPorts
}

func NetsConvertToStoreGraph(nets map[string]NetName) []*veryosys.NetName {
	netmaps := make(map[int]NetName)

	var graphNets []*veryosys.NetName
	var count int
	for _, net := range nets {
		if net.HideName == 0 {
			continue
		}
		if len(net.Bits) > 1 {
			for numindex, num := range net.Bits {
				graphNets = append(graphNets, &veryosys.NetName{
					BitNum:    num,
					Name:      fmt.Sprintf("_%d_[%d]", count, numindex),
					WireGroup: fmt.Sprintf("_%d_", count),
				})
			}
		} else {
			graphNets = append(graphNets, &veryosys.NetName{
				BitNum:    net.Bits[0],
				Name:      fmt.Sprintf("_%d_", count),
				WireGroup: fmt.Sprintf("_%d_", count),
			})
		}
		count += 1
	}
	return graphNets
}

func CellsConvertToStoreGraph(cells map[string]Cell) []*veryosys.Cell {
	var graphCells []*veryosys.Cell
	for _, cell := range cells {
		if cell.HideName == 0 {
			continue
		}
		conns := make(map[int]struct{ Type string })
		for k, v := range cell.PortDirections {
			bitnum := cell.Connections[k]
			if len(bitnum) > 1 {
				err := fmt.Errorf("multibit not implement: %v", bitnum)
				panic(err.Error())
			}
			conns[bitnum[0]] = struct{ Type string }{
				Type: v,
			}
		}
		graphCells = append(graphCells, &veryosys.Cell{
			Type:        cell.Type,
			Connections: conns,
		})
	}
	return graphCells
}

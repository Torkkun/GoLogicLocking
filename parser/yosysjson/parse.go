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
			BitNum:    port.Bits[0], //PortsのBitはすべて同じと今は仮定する
			BitWidth:  len(port.Bits),
			Name:      name,
		})
	}
	return graphPorts
}

func NetsConvertToStoreGraph(nets map[string]NetName) []*veryosys.NetName {
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

func CellsConvertToStoreGraph(cells map[string]Cell) {

}

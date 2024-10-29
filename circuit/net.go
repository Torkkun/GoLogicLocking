package circuit

import (
	"context"
	"errors"
	"goll/graph/veryosys"
	"goll/parser/yosysjson"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func NetCreate(tx neo4j.ExplicitTransaction, module *yosysjson.Module, ctx context.Context, tmpmap map[int]*Mapv) (map[int]*Mapv, error) {
	gname := yosysjson.NetsConvertToStoreGraph(module.NetName)
	// ノードカウント変数
	var nodecount int
	// 差分の値
	diffnum := len(tmpmap)

	for _, net := range gname {
		newnet := new(veryosys.DBNetName)
		newnet.Netname = net.Netname
		newnet.Attributes.Src = net.Attributes.Src
		nodecount += len(net.Bits)
		for _, bitnum := range net.Bits {
			newnet.BitNum = bitnum
			elementid, err := veryosys.CreateWireNodeTx(tx, ctx, newnet)
			if err != nil {
				return nil, err
			}
			tmpmap[bitnum] = &Mapv{
				Type:      "Wire",
				ElementId: elementid,
			}
		}
	}
	if len(tmpmap)-diffnum != nodecount {
		err := errors.New("net to Graph Error")
		return nil, err
	}
	log.Println("Create Wire Node OK")
	return tmpmap, nil
}

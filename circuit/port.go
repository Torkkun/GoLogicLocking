package circuit

import (
	"context"
	"errors"
	"goll/graph/veryosys"
	"goll/parser/yosysjson"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// json情報からグラフに変換
func PortCreate(tx neo4j.ExplicitTransaction, module *yosysjson.Module, ctx context.Context, tmpmap map[int]*Mapv) (map[int]*Mapv, error) {
	gport := yosysjson.PortsConvertToStoreGraph(module.Ports)
	for _, io := range gport {
		elementid, err := veryosys.CreateInOutNodeTx(tx, ctx, io)
		if err != nil {
			return nil, err
		}
		tmpmap[io.BitNum] = &Mapv{
			Type:      "IO",
			ElementId: elementid,
		}
	}
	if len(tmpmap) != len(gport) {
		err := errors.New("port to Graph Error")
		return nil, err
	}
	log.Println("Create IOPort Node OK")
	return tmpmap, nil
}

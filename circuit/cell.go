package circuit

import (
	"context"
	"goll/graph/veryosys"
	"goll/parser/yosysjson"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateCellAndConnectNodes(tx neo4j.ExplicitTransaction, module *yosysjson.Module, ctx context.Context, tmpmap map[int]*Mapv) error {
	gcell, gconns := yosysjson.CellsConvertToStoreGraph(module.Cells)

	cellnodes, err := CreateCell(tx, ctx, gcell)
	if err != nil {
		return err
	}
	if err = CellConnect(tx, ctx, gconns, cellnodes, tmpmap); err != nil {
		return err
	}
	return nil
}

func CreateCell(tx neo4j.ExplicitTransaction, ctx context.Context, gcell []*veryosys.Cell) (map[string]string, error) {
	cellnodes := make(map[string]string) // K AttrSrc, V elementId
	for _, cell := range gcell {
		elementid, err := veryosys.CreateCellNodeTx(tx, ctx, cell)
		if err != nil {
			return nil, err
		}
		cellnodes[cell.Attributes.Src] = elementid
	}
	log.Println("Create CellPort Node OK")
	return cellnodes, nil
}

func CellConnect(tx neo4j.ExplicitTransaction, ctx context.Context, gconns veryosys.Connections, cellnodes map[string]string, tmpmap map[int]*Mapv) error {
	for src, conn := range gconns {
		elementid := cellnodes[src]
		for _, v := range conn {
			p := tmpmap[v.BitNum]
			if v.Type == "input" {
				err := veryosys.CellConnectionTx(tx, ctx, &veryosys.ConnectionPair{
					Predecessor: veryosys.Node{
						Type:      "Cell",
						ElementId: elementid,
					},
					Successor: veryosys.Node{
						Type:      p.Type,
						ElementId: p.ElementId,
					},
				})
				if err != nil {
					log.Fatalf("Create Connection Node Error: %v", err.Error())
				}
			} else if v.Type == "output" {
				// Netで作成したものとCellを接続する
				err := veryosys.CellConnectionTx(tx, ctx, &veryosys.ConnectionPair{
					Predecessor: veryosys.Node{
						Type:      p.Type,
						ElementId: p.ElementId,
					},
					Successor: veryosys.Node{
						Type:      "Cell",
						ElementId: elementid,
					},
				})
				if err != nil {
					return err
				}
			}
		}
	}
	log.Println("Connect Node OK")
	return nil
}

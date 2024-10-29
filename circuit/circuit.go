package circuit

import (
	"context"
	"goll/graph/veryosys"
	"goll/parser/yosysjson"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Mapv struct {
	Type      string
	ElementId string
}

func TranslationCircuit(ctx context.Context, driver neo4j.DriverWithContext, dbname string, module *yosysjson.Module) error {
	// session 作成
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: dbname})
	defer session.Close(ctx)

	tmpmap := make(map[int]*Mapv)

	_, err := veryosys.Transaction(ctx, session, dbname, func(tx neo4j.ExplicitTransaction) (interface{}, error) {
		tmpmap, err := PortCreate(tx, module, ctx, tmpmap)
		if err != nil {
			return nil, err
		}
		tmpmap, err = NetCreate(tx, module, ctx, tmpmap)
		if err != nil {
			return nil, err
		}
		if err = CreateCellAndConnectNodes(tx, module, ctx, tmpmap); err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type IONode struct {
	Type string
	Name string
}

func (io *IONode) CreateInOutNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:IO {type: $type, name:$name})`,
		map[string]any{
			"type": io.Type,
			"name": io.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return err
	}
	return nil
}

type LogicGateNode struct {
	GateType string
	At       int
}

func (gate *LogicGateNode) CreateLogicGateNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Gate {type: $type, at: $at})`,
		map[string]any{
			"type": gate.GateType,
			"at":   gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateLogicGateNode Error:%v", err)
		return err
	}
	return nil
}

type WireNode struct {
	Name string
}

func (wire *WireNode) CreateWireNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Wire {name: $name})`,
		map[string]any{
			"name": wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateWireNode Error:%v", err)
		return err
	}
	return nil
}

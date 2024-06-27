package verpyverilog

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type IONode struct {
	Type string
	Name string
}

func (io *IONode) CreateInOutNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (io:IO {type: $type, name:$name})
		RETURN io`,
		map[string]any{
			"type": io.Type,
			"name": io.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return "", err
	}
	if len(result.Records) != 1 {
		err = fmt.Errorf("to match node:%v", result.Records)
		return "", err
	}
	record := result.Records[0]
	resio, ok := record.Get("io")
	if !ok {
		err = fmt.Errorf("NotfoundCreated IO Node")
		return "", err
	}
	n := resio.(neo4j.Node)
	return n.GetElementId(), nil
}

type LogicGateNode struct {
	GateType string
}

func (gate *LogicGateNode) CreateLogicGateNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (g:Gate {type: $type})
		RETURN g`,
		map[string]any{
			"type": gate.GateType,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateLogicGateNode Error:%v", err)
		return "", err
	}
	if len(result.Records) != 1 {
		err = fmt.Errorf("to match node:%v", result.Records)
		return "", err
	}
	record := result.Records[0]
	resg, ok := record.Get("g")
	if !ok {
		err = fmt.Errorf("NotfoundCreated Logic Node")
		return "", err
	}
	n := resg.(neo4j.Node)
	return n.GetElementId(), nil
}

type LockGateNode struct {
	GateType string
	LockType string
}

func (lg *LockGateNode) CreateLockingGateNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (g:LLGate {type: $type, ll: $ll})
		RETURN g`,
		map[string]any{
			"type": lg.GateType,
			"ll":   lg.LockType,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateLockingGateNode Error:%v", err)
		return "", err
	}
	if len(result.Records) != 1 {
		err = fmt.Errorf("to match node:%v", result.Records)
		return "", err
	}
	record := result.Records[0]
	reslg, ok := record.Get("g")
	if !ok {
		err = fmt.Errorf("NotfoundCreated Locking Node")
		return "", err
	}
	n := reslg.(neo4j.Node)
	return n.GetElementId(), nil
}

type WireNode struct {
	Name string
}

func (wire *WireNode) CreateWireNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (w:Wire {name: $name})
		RETURN w`,
		map[string]any{
			"name": wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateWireNode Error:%v", err)
		return "", err
	}
	if len(result.Records) != 1 {
		return "", fmt.Errorf("too much node: %v", result.Records)
	}
	resw, ok := result.Records[0].Get("w")
	if !ok {
		err = fmt.Errorf("NotfoundCreated Wire Node")
		return "", err
	}
	n := resw.(neo4j.Node)
	return n.GetElementId(), nil
}

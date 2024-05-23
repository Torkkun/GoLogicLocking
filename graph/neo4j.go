package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type graphDB struct {
	Driver neo4j.DriverWithContext
}

func NewDriver() *graphDB {
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	return &graphDB{Driver: driver}
}

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

// 後で良い方法を考える OGM（object graph mapping）
// IO(IN) <- Gate

type GateIO struct {
	Gate LogicGateNode
	Io   IONode
	At   int
}

func (gi *GateIO) GatetoIO(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (io:IO {type: $io_type, name: $io_name}), (g:Gate {type: $g_type, at: $g_at})
		MERGE (io)<-[:LGtoIO]-(g)
		`,
		map[string]any{
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
			"g_type":  gi.Gate.GateType,
			"g_at":    gi.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE GatetoIO Error:%v", err)
		return err
	}
	return nil
}

// Gate <- IO(OUT)
func (gi *GateIO) IOtoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at}), (io:IO {type: $io_type, name: $io_name})
		MERGE (g)<-[:IOtoLG]-(io)`,
		map[string]any{
			"g_type":  gi.Gate.GateType,
			"g_at":    gi.Gate.At,
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE IOtoGate Error:%v", err)
		return err
	}
	return nil
}

type GateWire struct {
	Gate LogicGateNode
	Wire WireNode
	At   int
}

// Gate <- Wire
func (gw *GateWire) WiretoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at}), (w:Wire {name: $w_name})
		MERGE (g)<-[:WiretoLG]-(w)`,
		map[string]any{
			"g_type": gw.Gate.GateType,
			"g_at":   gw.Gate.At,
			"w_name": gw.Wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE WiretoGate Error:%v", err)
		return err
	}
	return nil
}

// Wire <- Gate
func (gw *GateWire) GatetoWire(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (w:Wire {name: $w_name}), (g:Gate {type: $g_type, at: $g_at})
		MERGE (w)<-[:LGtoWire]-(g)`,
		map[string]any{
			"w_name": gw.Wire.Name,
			"g_type": gw.Gate.GateType,
			"g_at":   gw.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return err
	}
	return nil
}

// 複数のユーザーデータベースは作れないので全部消すしかないため用意
func DBtableAllClear(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		``,
		map[string]any{},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DBtableAllClear Error:%v", err)
		return err
	}
	return nil
}

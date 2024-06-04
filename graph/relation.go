package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetLGtoIORelation() {
	//MATCH (:IO {name: $io_name, type: $io_type})<-[r:LGtoIO]-(:Gate) RETURN r
}

// Create Relationship
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
		err = fmt.Errorf("MERGE GatetoWire Error:%v", err)
		return err
	}
	return nil
}

// Delete Relationship

func (gi *GateIO) DeleteRelationGatetoIO(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (io:IO {type: $io_type, name: $io_name})<-[r:LGtoIO]-(g:Gate {type: $g_type, at: $g_at})
		DELETE r`,
		map[string]any{
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
			"g_type":  gi.Gate.GateType,
			"g_at":    gi.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DELETE GatetoIO Relation Error:%v", err)
		return err
	}
	return nil
}

func (gi *GateIO) DeleteRelationIOtoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at})<-[r:IOtoLG]-(io:IO {type: $io_type, name: $io_name})
		DELETE r`,
		map[string]any{
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
			"g_type":  gi.Gate.GateType,
			"g_at":    gi.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DELETE IOtoGate Relation Error:%v", err)
		return err
	}
	return nil
}

func (gw *GateWire) DeleteRelationWiretoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at})<-[r:WiretoLG]-(w:Wire {name: $w_name})
		DELETE r`,
		map[string]any{
			"g_type": gw.Gate.GateType,
			"g_at":   gw.Gate.At,
			"w_name": gw.Wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("DELETE WiretoGate Relation Error:%v", err)
		return err
	}
	return nil
}

func (gw *GateWire) DeleteRelationGatetoWire(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (w:Wire {name: $w_name})<-[r:LGtoWire]-(g:Gate {type: $g_type, at: $g_at})
		DELETE r `,
		map[string]any{
			"w_name": gw.Wire.Name,
			"g_type": gw.Gate.GateType,
			"g_at":   gw.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("DELETE GatetoWire Relation Error:%v", err)
		return err
	}
	return nil
}

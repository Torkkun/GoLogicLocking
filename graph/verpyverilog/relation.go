package verpyverilog

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Relation struct {
	Identity           int64
	Type               string
	Properties         interface{}
	ElementId          string
	StartNodeElementId string
	EndNodeElementId   string
}

// Create Relationship
// 後で良い方法を考える OGM（object graph mapping）
// IO(IN) <- Gate

type GateIO struct {
	Gate *LogicGateNode
	Io   *IONode
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

type LockGateIO struct {
	Gate *LockGateNode
	Io   *IONode
}

func (lgi *LockGateIO) LLGatetoIO(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (io:IO {type: $io_type, name: $io_name}), (g:LLGate {type: $g_type, locktype: $ll_type, name: $name})
		MERGE (io)<-[:LLGtoIO]-(g)
		`,
		map[string]any{
			"io_type": lgi.Io.Type,
			"io_name": lgi.Io.Name,
			"g_type":  lgi.Gate.GateType,
			"ll_type": lgi.Gate.LockType,
			"name":    lgi.Gate.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
}

func (lgi *LockGateIO) LLGatetoIOByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (io:IO {type: $io_type, name: $io_name}), (lg:LLGate)
		WHERE elementId(lg)=$element_id
		MERGE (io)<-[:LLGtoIO]-(lg)
		`,
		map[string]any{
			"io_type":    lgi.Io.Type,
			"io_name":    lgi.Io.Name,
			"element_id": lgi.Gate.ElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
}

// IO(IN) <- Gateのリレーションを取得
func (io *IONode) GetLGtoIORelation(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (:IO {name: $io_name, type: $io_type})<-[r:LGtoIO]-(:Gate) RETURN r`,
		map[string]any{
			"io_name": io.Name,
			"io_type": io.Type,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH LGtoIO Relation Error:%v", err)
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

// LockingGate <- IO(OUT)
func (lgi *LockGateIO) IOtoLLGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (lg:LLGate), (io:IO {type: $io_type, name: $io_name})
		WHERE elementId(lg)=$element_id
		MERGE (lg)<-[:IOtoLLG]-(io)`,
		map[string]any{
			"io_type":    lgi.Io.Type,
			"io_name":    lgi.Io.Name,
			"element_id": lgi.Gate.ElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE IOtoLLGate Error:%v", err)
		return err
	}
	return nil
}

type GetNeo4JIoAndRelation struct {
	Neo4JIO  *GetNeo4JIONode
	Relation *Relation
}

// Gate <- IO(OUT)のリレーション
func (lg *LogicGateNode) GetIOtoGateRelation(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*GetNeo4JIoAndRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`
		MATCH (:Gate {type: $g_type, at: $g_at})<-[r:IOtoLG]-(io:IO)
		RETURN r,io
		`,
		map[string]any{
			"g_type": lg.GateType,
			"g_at":   lg.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH IOtoLG Relation Error:%v", err)
		return nil, err
	}
	ioandrelation := make(map[string]*GetNeo4JIoAndRelation)
	for _, record := range result.Records {
		io, ok := record.Get("io")
		if !ok {
			err = fmt.Errorf("GetPredecessor IO Node Error")
			return nil, err
		}
		tmpio := io.(neo4j.Node)
		r, ok := record.Get("r")
		if !ok {
			err = fmt.Errorf("GetIO to Gate Relation Error")
			return nil, err
		}
		tmpr := r.(neo4j.Relationship)
		rElementId := tmpr.GetElementId()
		ioandrelation[rElementId] = &GetNeo4JIoAndRelation{
			Neo4JIO: &GetNeo4JIONode{
				ION: &IONode{
					Type: tmpio.Props["type"].(string),
					Name: tmpio.Props["name"].(string),
				},
				ElementId: tmpio.GetElementId(),
				Id:        tmpio.GetId(),
			},
			Relation: &Relation{
				Identity:           tmpr.GetId(),
				ElementId:          rElementId,
				StartNodeElementId: tmpr.StartElementId,
				EndNodeElementId:   tmpr.EndElementId,
			},
		}
	}
	return ioandrelation, nil
}

func GetIOtoGateRelationByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) ([]*GetNeo4JIoAndRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`
		MATCH (g:Gate)<-[r:IOtoLG]-(io:IO)
		WHERE elementId(g)=$element_id
		RETURN r,io
		`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH IOtoLG Relation Error:%v", err)
		return nil, err
	}
	ioandrelation := []*GetNeo4JIoAndRelation{}
	for _, record := range result.Records {
		io, ok := record.Get("io")
		if !ok {
			err = fmt.Errorf("GetPredecessor IO Node Error")
			return nil, err
		}
		tmpio := io.(neo4j.Node)
		r, ok := record.Get("r")
		if !ok {
			err = fmt.Errorf("GetIO to Gate Relation Error")
			return nil, err
		}
		tmpr := r.(neo4j.Relationship)
		ioandrelation = append(ioandrelation, &GetNeo4JIoAndRelation{
			Neo4JIO: &GetNeo4JIONode{
				ION: &IONode{
					Type: tmpio.Props["type"].(string),
					Name: tmpio.Props["name"].(string),
				},
				ElementId: tmpio.GetElementId(),
				Id:        tmpio.GetId(),
			},
			Relation: &Relation{
				Identity:           tmpr.GetId(),
				ElementId:          tmpr.GetElementId(),
				StartNodeElementId: tmpr.StartElementId,
				EndNodeElementId:   tmpr.EndElementId,
			},
		})
	}

	return ioandrelation, nil
}

type GateWire struct {
	Gate *LogicGateNode
	Wire *WireNode
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

type LockGateWire struct {
	Gate *LockGateNode
	Wire *WireNode
}

// LGate <- Wire
func (lgw *LockGateWire) WiretoLLGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (lg:LLGate), (w:Wire {name: $w_name})
		WHERE elementId(lg)=$element_id
		MERGE (lg)<-[:WiretoLLG]-(w)
		`,
		map[string]any{
			"w_name":     lgw.Wire.Name,
			"element_id": lgw.Gate.ElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
}

type GetNeo4JWireAndRelation struct {
	Neo4JWire *GetNeo4JWireNode
	Relation  *Relation
}

// Gate <- Wireのリレーションを取得
func (g *LogicGateNode) GetWiretoGateRelation(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*GetNeo4JWireAndRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (:Gate {type: $g_type, at: $g_at})<-[r:WiretoLG]-(w:Wire) RETURN r,w`,
		map[string]any{
			"g_type": g.GateType,
			"g_at":   g.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH WiretoGate Relation Error:%v", err)
		return nil, err
	}
	wireandrelation := make(map[string]*GetNeo4JWireAndRelation)
	for _, record := range result.Records {
		w, ok := record.Get("w")
		if !ok {
			err = fmt.Errorf("GetPredecessor Wire Node Error")
			return nil, err
		}
		tmpw := w.(neo4j.Node)
		r, ok := record.Get("r")
		if !ok {
			err = fmt.Errorf("GetWire to Gate Relation Error")
			return nil, err
		}
		tmpr := r.(neo4j.Relationship)
		rElementId := tmpr.GetElementId()
		wireandrelation[rElementId] = &GetNeo4JWireAndRelation{
			Neo4JWire: &GetNeo4JWireNode{
				WN: &WireNode{
					Name: tmpw.Props["name"].(string),
				},
				ElementId: tmpw.GetElementId(),
				Id:        tmpw.GetId(),
			},
			Relation: &Relation{
				Identity:           tmpr.GetId(),
				ElementId:          tmpr.GetElementId(),
				StartNodeElementId: tmpr.StartElementId,
				EndNodeElementId:   tmpr.EndElementId,
			},
		}

	}
	return wireandrelation, nil
}

func GetWiretoGateRelationByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) ([]*GetNeo4JWireAndRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`
		MATCH (g:Gate)<-[r:WiretoLG]-(w:Wire)
		WHERE elementId(g)=$element_id
		RETURN r,w
		`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH WiretoGate Relation Error:%v", err)
		return nil, err
	}
	wireandrelation := []*GetNeo4JWireAndRelation{}
	for _, record := range result.Records {
		w, ok := record.Get("w")
		if !ok {
			err = fmt.Errorf("GetPredecessor Wire Node Error")
			return nil, err
		}
		tmpw := w.(neo4j.Node)
		r, ok := record.Get("r")
		if !ok {
			err = fmt.Errorf("GetWire to Gate Relation Error")
			return nil, err
		}
		tmpr := r.(neo4j.Relationship)
		wireandrelation = append(wireandrelation, &GetNeo4JWireAndRelation{
			Neo4JWire: &GetNeo4JWireNode{
				WN: &WireNode{
					Name: tmpw.Props["name"].(string),
				},
				ElementId: tmpw.GetElementId(),
				Id:        tmpw.GetId(),
			},
			Relation: &Relation{
				Identity:           tmpr.GetId(),
				ElementId:          tmpr.GetElementId(),
				StartNodeElementId: tmpr.StartElementId,
				EndNodeElementId:   tmpr.EndElementId,
			},
		})

	}
	return wireandrelation, nil
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

type LLGateGate struct {
	LLGN *LockGateNode
	LGN  *LogicGateNode
}

// Gate <- LLGate
func (llgg *LLGateGate) LLGatetoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at}), (lg:LLGate {type: $lg_type, locktype: $ll_type, name: $name})
		MERGE (g)<-[:LLGtoG]-(lg)`,
		map[string]any{
			"g_type":  llgg.LGN.GateType,
			"g_at":    llgg.LGN.At,
			"lg_type": llgg.LLGN.GateType,
			"ll_type": llgg.LLGN.LockType,
			"name":    llgg.LLGN.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoGate Error:%v", err)
		return err
	}
	return nil
}

func (llgg *LLGateGate) LLGatetoGateElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at}), (lg:LLGate)
		WHERE elementId(lg)=$element_id
		MERGE (g)<-[:LLGtoG]-(lg)`,
		map[string]any{
			"g_type":     llgg.LGN.GateType,
			"g_at":       llgg.LGN.At,
			"element_id": llgg.LLGN.ElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoGate Error:%v", err)
		return err
	}
	return nil
}

package verpyverilog

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GetNeo4JLogicNode struct {
	LGN       *LogicGateNode
	Id        int64
	ElementId string
}

func GetAllLogicNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*GetNeo4JLogicNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(g:Gate) RETURN g`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllLogicNodes Error:%v", err)
		return nil, err
	}
	lgnodes := make(map[string]*GetNeo4JLogicNode)
	for _, record := range result.Records {
		io, ok := record.Get("g")
		if !ok {
			err = fmt.Errorf("GetAllLogicNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		id := tmp.GetId()
		lgnodes[elementId] = &GetNeo4JLogicNode{
			LGN: &LogicGateNode{
				GateType: tmp.Props["type"].(string),
				At:       int(tmp.Props["at"].(int64)),
			},
			ElementId: elementId,
			Id:        id,
		}
	}
	return lgnodes, nil
}

type GetNeo4JIONode struct {
	ION       *IONode
	Id        int64
	ElementId string
}

func GetAllIONode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*GetNeo4JIONode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(io:IO) RETURN io`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllIONodes Error:%v", err)
		return nil, err
	}
	ionodes := make(map[string]*GetNeo4JIONode)
	for _, record := range result.Records {
		io, ok := record.Get("io")
		if !ok {
			err = fmt.Errorf("GetAllIONode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		id := tmp.GetId()
		ionodes[elementId] = &GetNeo4JIONode{
			ION: &IONode{
				Type: tmp.Props["type"].(string),
				Name: tmp.Props["name"].(string),
			},
			ElementId: elementId,
			Id:        id,
		}
	}
	return ionodes, nil
}

type GetNeo4JWireNode struct {
	WN        *WireNode
	Id        int64
	ElementId string
}

func GetAllWireNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*GetNeo4JWireNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(w:Wire) RETURN w`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllNodes Error:%v", err)
		return nil, err
	}
	wnodes := make(map[string]*GetNeo4JWireNode)
	for _, record := range result.Records {
		io, ok := record.Get("w")
		if !ok {
			err = fmt.Errorf("GetAllWireNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		id := tmp.GetId()
		wnodes[elementId] = &GetNeo4JWireNode{
			WN: &WireNode{
				Name: tmp.Props["name"].(string),
			},
			ElementId: elementId,
			Id:        id,
		}
	}
	return wnodes, nil
}

type GetNeo4JLockGateNode struct {
	LockGateNode *LockGateNode
}

func GetAllLockGateNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*GetNeo4JLockGateNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(lg:LLGate) RETURN lg`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllLLGateNodes Error:%v", err)
		return nil, err
	}
	lgnodes := make(map[string]*GetNeo4JLockGateNode)
	for _, record := range result.Records {
		io, ok := record.Get("lg")
		if !ok {
			err = fmt.Errorf("GetAllLLGateNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		lgnodes[elementId] = &GetNeo4JLockGateNode{
			LockGateNode: &LockGateNode{
				Name:      tmp.Props["name"].(string),
				LockType:  tmp.Props["ll"].(string),
				GateType:  tmp.Props["type"].(string),
				ElementId: elementId,
			},
		}
	}
	return lgnodes, nil
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

type PredecessorNode struct {
	RelationElementID string
	NodeElementID     string
}

// A <- (B) ,Get "B" Node
func GetAllPredecessors(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) ([]*PredecessorNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:%)<-[r:%]-(n:%)
		WHERE elementId(g)=$element_id
		RETURN r,n`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH Get All Predecessors Error:%v", err)
		return nil, err
	}
	predecessors := []*PredecessorNode{}
	for _, record := range result.Records {
		node, ok := record.Get("n")
		if !ok {
			err = fmt.Errorf("GetPredecessor Node Error")
			return nil, err
		}
		tmpn := node.(neo4j.Node)
		relation, ok := record.Get("r")
		if !ok {
			err = fmt.Errorf("GetPredecessor Relation Error")
			return nil, err
		}
		tmpr := relation.(neo4j.Relationship)
		predecessors = append(predecessors, &PredecessorNode{
			RelationElementID: tmpr.GetElementId(),
			NodeElementID:     tmpn.GetElementId(),
		})
	}
	return predecessors, nil
}

type SuccessorNode struct {
	RelationElementID string
	NodeElementID     string
	NodeName          string
}

// (C) <- A ,Get "C" Node
func GetAllSuccessorNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) ([]*SuccessorNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:%)-[r:%]->(n:%)
		WHERE elementId(g)=$element_id
		RETURN r,n`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH Get All Successor Error:%v", err)
		return nil, err
	}
	successors := []*SuccessorNode{}
	for _, record := range result.Records {
		node, ok := record.Get("n")
		if !ok {
			err = fmt.Errorf("GetSuccessor Node Error")
			return nil, err
		}
		tmpn := node.(neo4j.Node)
		relation, ok := record.Get("r")
		if !ok {
			err = fmt.Errorf("GetSuccessor Relation Error")
			return nil, err
		}
		tmpr := relation.(neo4j.Relationship)
		successors = append(successors, &SuccessorNode{
			RelationElementID: tmpr.GetElementId(),
			NodeElementID:     tmpn.GetElementId(),
			NodeName:          tmpn.Props["name"].(string),
		})
	}
	return successors, nil
}

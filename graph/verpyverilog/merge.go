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

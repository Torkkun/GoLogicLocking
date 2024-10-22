package veryosys

import (
	"context"
	"errors"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func GetInNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*Port, error) {
	query := `
	MATCH(io:IO WHERE io.direction = 'input')
	RETURN io`
	return GetPortNode(ctx, driver, dbname, query)
}
func GetOutNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*Port, error) {
	query := `
	MATCH(io:IO WHERE io.direction = 'output')
	RETURN io`
	return GetPortNode(ctx, driver, dbname, query)
}

func GetPortNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, query string) (map[string]*Port, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		query,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetInNode Error:%v", err)
		return nil, err
	}
	ionodes := make(map[string]*Port)
	for _, record := range result.Records {
		io, ok := record.Get("io")
		if !ok {
			err = fmt.Errorf("GetInNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		ionodes[elementId] = &Port{
			Direction: tmp.Props["direction"].(string),
			Name:      tmp.Props["name"].(string),
			BitNum:    tmp.Props["bitnum"].(int),
			BitWidth:  tmp.Props["bitwidth"].(int),
		}
	}
	return ionodes, nil
}

func GetWireNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*DBNetName, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(w:Wire)
		RETURN w`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetNetNode Error:%v", err)
		return nil, err
	}
	netnodes := make(map[string]*DBNetName)
	for _, record := range result.Records {
		io, ok := record.Get("io")
		if !ok {
			err = fmt.Errorf("GetNetNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		netnodes[elementId] = &DBNetName{
			BitNum:  tmp.Props["bitnum"].(int),
			Netname: tmp.Props["netname"].(string),
			Attributes: struct{ Src string }{
				tmp.Props["attrsrc"].(string),
			},
		}
	}
	return netnodes, nil
}

func GetCellNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (map[string]*Cell, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(cell:Cell)
		RETURN wcell`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetNetNode Error:%v", err)
		return nil, err
	}
	cellnodes := make(map[string]*Cell)
	for _, record := range result.Records {
		io, ok := record.Get("cell")
		if !ok {
			err = fmt.Errorf("GetCellNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		elementId := tmp.GetElementId()
		cellnodes[elementId] = &Cell{
			Type: tmp.Props["type"].(string),
			Attributes: struct{ Src string }{
				tmp.Props["attrsrc"].(string),
			},
		}
	}
	return cellnodes, nil
}

type GetPreNodeAndRelation struct {
	Node     []neo4j.Node
	Relation []neo4j.Relationship
}

func GetPredecessorNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string, sucElementId string) (*GetPreNodeAndRelation, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver, `
	MATCH (pre)-[r]->(suc)
	WHERE elementId(suc)=$suc_element_id
	RETURN r, pre`,
		map[string]any{
			"suc_element_id": sucElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH PredecessorNode Error:%v", err)
		return nil, err
	}
	var nodes []neo4j.Node
	var relation []neo4j.Relationship
	for _, record := range result.Records {
		// convert to node
		preres, ok := record.Get("pre")
		if !ok {
			err = errors.New("get predecessor node record error")
			return nil, err
		}
		prenode := preres.(neo4j.Node)
		nodes = append(nodes, prenode)
		// convert to relationship
		rres, ok := record.Get("r")
		if !ok {
			err = errors.New("get Relation Record Error")
			return nil, err
		}
		r := rres.(neo4j.Relationship)
		relation = append(relation, r)
	}
	return &GetPreNodeAndRelation{
		Node:     nodes,
		Relation: relation,
	}, nil
}

type GetRelationships struct {
	Pre      []neo4j.Node
	Suc      []neo4j.Node
	Relation []neo4j.Relationship
}

// WiretoCellだけ条件に含めない
func GetRelationshipsAndPairNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (*GetRelationships, error) {
	// WiretoCellだけとりのぞく
	result, err := neo4j.ExecuteQuery(ctx, driver, `
	MATCH (pre)-[r]->(suc)
	WHERE NOT type(r) = 'WiretoCell'
	RETURN pre, r, suc`,
		map[string]any{},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MATCH Relationship Error:%v", err)
		return nil, err
	}
	var prenodes []neo4j.Node
	var sucnodes []neo4j.Node
	var relation []neo4j.Relationship
	for _, record := range result.Records {
		// convert to pre node
		preres, ok := record.Get("pre")
		if !ok {
			err = errors.New("get predecessor node record error")
			return nil, err
		}
		prenode := preres.(neo4j.Node)
		prenodes = append(prenodes, prenode)
		// convert to suc node
		sucres, ok := record.Get("suc")
		if !ok {
			err = errors.New("get successor node record error")
			return nil, err
		}
		sucnode := sucres.(neo4j.Node)
		sucnodes = append(sucnodes, sucnode)
		// convert to relationship
		rres, ok := record.Get("r")
		if !ok {
			err = errors.New("get Relation Record Error")
			return nil, err
		}
		r := rres.(neo4j.Relationship)
		relation = append(relation, r)
	}
	return &GetRelationships{
		Pre:      prenodes,
		Suc:      sucnodes,
		Relation: relation,
	}, nil
}

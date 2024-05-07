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

func MatchRelationship() {

}

type IONode struct {
	Type string
	Name string
}

func (io *IONode) CreateInOutNode(ctx context.Context, driver neo4j.DriverWithContext) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MERGE (io:IO {type: $type, name:$name}) RETURN io`,
		map[string]any{
			"type": io.Type,
			"name": io.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j")) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return err
	}
	return nil
}

type GateNode struct {
	GateType string
	At       int
}

func (gate *GateNode) CreateGateNode(ctx context.Context, driver neo4j.DriverWithContext) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MERGE (g:Gate {type: $type, at: $at}) RETURN g`,
		map[string]any{
			"type": gate.GateType,
			"at":   gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase("neo4j")) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return err
	}
	return nil
}

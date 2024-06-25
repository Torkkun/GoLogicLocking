package verpyverilog

import (
	"context"
	"fmt"
	"goll/utils"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestRelationIOtoLG(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)
	lg := LogicGateNode{GateType: "OR", At: 23}
	lg.GetIOtoGateRelation(ctx, driver, "neo4j")
}

func TestRelationLGtoW(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)
}

func TestGetIOtoGateRelation(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	g := new(LogicGateNode)
	g.At = 25
	g.GateType = string(utils.Or)
	relation, err := g.GetIOtoGateRelation(ctx, driver, "neo4j")
	if err != nil {
		panic(err)
	}
	for _, r := range relation {
		fmt.Println(r.Neo4JIO)
		fmt.Println(r.Relation)
	}
}

func TestGetWiretoGateRelation(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	g := new(LogicGateNode)
	g.At = 23
	g.GateType = string(utils.Or)
	relation, err := g.GetWiretoGateRelation(ctx, driver, "neo4j")
	if err != nil {
		panic(err)
	}
	for _, r := range relation {
		fmt.Println(r.Neo4JWire)
		fmt.Println(r.Relation)
	}
}

func TestGetWtoGRelationByElementId(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	elementId := "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:47"
	relation, err := GetWiretoGateRelationByElementId(ctx, driver, "neo4j", elementId)
	if err != nil {
		panic(err)
	}
	for _, r := range relation {
		fmt.Println(r.Neo4JWire)
		fmt.Println(r.Relation)
	}
}

func TestGetIOtoGRelationByElementId(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	elementId := "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:42"
	relation, err := GetIOtoGateRelationByElementId(ctx, driver, "neo4j", elementId)
	if err != nil {
		panic(err)
	}
	for _, r := range relation {
		fmt.Println(r.Neo4JIO)
		fmt.Println(r.Relation)
	}
}

func TestGetGatePredecessorNodes(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	elementId := "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:47"
	relation1, err := GetWiretoGateRelationByElementId(ctx, driver, "neo4j", elementId)
	if err != nil {
		panic(err)
	}
	relation2, err := GetIOtoGateRelationByElementId(ctx, driver, "neo4j", elementId)
	if err != nil {
		panic(err)
	}
	fmt.Println(relation1)
	fmt.Println(relation2)
	if relation1 != nil {
		fmt.Println("OK")
	}
}

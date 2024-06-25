package graph

import (
	"context"
	"fmt"
	"goll/graph/verpyverilog"
	"log"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestConn(t *testing.T) {
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
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}
}

func TestCountGate(t *testing.T) {
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
	num, err := CountGate(ctx, driver, "neo4j")
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
}

func TestCountOUT(t *testing.T) {
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
	num, err := CountOUT(ctx, driver, "neo4j")
	if err != nil {
		panic(err)
	}
	fmt.Println(num)
}

func TestGetLG(t *testing.T) {
	ctx := context.Background()
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		log.Fatalln(err)
	}
	defer driver.Close(ctx)
	lgn, err := verpyverilog.GetAllLogicNodes(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	for _, lg := range lgn {
		fmt.Println(*lg.LGN, lg.ElementId, lg.Id)
	}
}

func TestGetIO(t *testing.T) {
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
	ion, err := verpyverilog.GetAllIONode(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	for _, io := range ion {
		fmt.Println(*io.ION, io.ElementId, io.Id)
	}
}

func TestGetWire(t *testing.T) {
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
	wn, err := verpyverilog.GetAllWireNodes(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	for _, w := range wn {
		fmt.Println(*w.WN, w.ElementId, w.Id)
	}
}

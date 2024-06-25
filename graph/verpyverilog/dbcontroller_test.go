package verpyverilog

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestGetAllNode(t *testing.T) {
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
	all, err := GetAllNodes(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	for _, val := range all.Lgs {
		fmt.Println(*val.LGN)
	}
	for _, val := range all.Ws {
		fmt.Println(*val.WN)
	}
}

func TestSeparateIOType(t *testing.T) {
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
	ionodes, err := GetAllIONode(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	separatedio, err := SeparateIOType(ionodes)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("IN nodes")
	for _, ios := range separatedio.In {
		fmt.Println(*ios.ION)
	}
	fmt.Println("OUT nodes")
	for _, ios := range separatedio.Out {
		fmt.Println(*ios.ION)
	}
}

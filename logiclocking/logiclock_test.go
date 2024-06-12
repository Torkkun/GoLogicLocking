package logiclocking

import (
	"context"
	"fmt"
	"goll/graph"
	"goll/utils"
	"log"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func TestXOR(t *testing.T) {
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

	all, err := graph.GetAllNodes(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	var idlist []Id
	for _, v := range all.Lgs {
		idlist = append(idlist, Id{
			ElementId: v.ElementId,
			Id:        v.Id,
		})
	}

	gates, err := utils.Sample(idlist, 4)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(gates)
	key := make(map[string]bool)
	for i, g := range gates {
		b, err := utils.Choice([]bool{true, false})
		if err != nil {
			panic(err)
		}
		keystring := fmt.Sprintf("key_%v", i)
		key[keystring] = b
		pres, err := GetGatePredecessorNodes(ctx, driver, "neo4j", g)
		if err != nil {
			panic(err)
		}
		fmt.Println(i)
		if len(pres.IOaR) != 0 {
			fmt.Println(pres.IOaR)
		}
		if len(pres.WaR) != 0 {
			fmt.Println(pres.WaR)
		}
	}

}

func TestXorLock(t *testing.T) {
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
	if err = XorLock(ctx, driver, "neo4j", 2); err != nil {
		log.Fatalln(err)
	}
}

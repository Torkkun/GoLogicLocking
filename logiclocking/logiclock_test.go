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
	type Id struct {
		ElementId string
		Id        int64
	}
	var idlist []Id
	for _, v := range all.Lgs {
		idlist = append(idlist, Id{
			ElementId: v.ElementId,
			Id:        v.Id,
		})
	}
	for _, v := range all.Ws {
		idlist = append(idlist, Id{
			ElementId: v.ElementId,
			Id:        v.Id,
		})
	}
	for _, v := range all.Ios.In {
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
}

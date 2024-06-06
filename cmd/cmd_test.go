package cmd

import (
	"context"
	"goll/graph"
	"goll/parser"
	"log"
	"testing"
)

func TestNewParseDB(t *testing.T) {
	parseresult := parser.NewParse("../testtxt/test.txt")
	// Connectionの設定は後で考える
	dbname := "neo4j"

	driver := graph.NewDriver()
	//driver := graph.SelectDriver("origin")
	ctx := context.Background()

	defer driver.Driver.Close(ctx)
	var err error
	for _, io := range parseresult.Declarations.IOPorts {
		neoio := graph.IONode{
			Type: string(io.Type),
			Name: io.Name,
		}
		if err = neoio.CreateInOutNode(ctx, driver.Driver, dbname); err != nil {
			log.Fatalln(err)
		}
	}
	for _, wire := range parseresult.Declarations.Wires {
		neowire := graph.WireNode{
			Name: wire.Name,
		}
		if err = neowire.CreateWireNode(ctx, driver.Driver, dbname); err != nil {
			log.Fatalln(err)
		}
	}
	for _, logicgate := range parseresult.LogicGates {
		gate := graph.LogicGateNode{
			GateType: logicgate.GateType,
			At:       logicgate.At,
		}
		if err = gate.CreateLogicGateNode(ctx, driver.Driver, dbname); err != nil {
			log.Fatalln(err)
		}
	}
	for at, relation := range parseresult.Nodes {
		// i1 <- lg, i2 <- lg, lg <- out
		// i1 <- lg
		if err := graph.LGtoIN(ctx, driver.Driver, dbname, relation.In1, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}

		// i2 <- lg
		if err := graph.LGtoIN(ctx, driver.Driver, dbname, relation.In2, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}

		// lg <- out
		if err := graph.OUTtoLG(ctx, driver.Driver, dbname, relation.Out, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}
	}
}

func TestDeleteAllDB(t *testing.T) {
	//driver := graph.SelectDriver("origin")
	//driver := graph.SelectDriver("copy")
	driver := graph.NewDriver()
	ctx := context.Background()

	defer driver.Driver.Close(ctx)
	var err error
	if err = graph.DBtableAllClear(ctx, driver.Driver, driver.DBname); err != nil {
		log.Fatal(err)
	}
}

func TestLl(t *testing.T) {

}

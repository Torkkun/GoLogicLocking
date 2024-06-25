package main

import (
	"context"
	"goll/graph"
	"goll/graph/verpyverilog"
	"goll/parser"
	"log"
)

//テスト用の実行プログラム

func fullAddrtestMain() {
	//testfilename := "test.txt"
	testfilename := "./testtxt/test.txt"
	parseresult := parser.NewParse(testfilename)
	dbname := "neo4j"
	driver := graph.NewDriver()
	ctx := context.Background()
	defer driver.Driver.Close(ctx)
	var err error
	for _, io := range parseresult.Declarations.IOPorts {
		neoio := verpyverilog.IONode{
			Type: string(io.Type),
			Name: io.Name,
		}
		if err = neoio.CreateInOutNode(ctx, driver.Driver, dbname); err != nil {
			log.Fatalln(err)
		}
	}
	for _, wire := range parseresult.Declarations.Wires {
		neowire := verpyverilog.WireNode{
			Name: wire.Name,
		}
		if err = neowire.CreateWireNode(ctx, driver.Driver, dbname); err != nil {
			log.Fatalln(err)
		}
	}
	for _, logicgate := range parseresult.LogicGates {
		gate := verpyverilog.LogicGateNode{
			GateType: string(logicgate.GateType),
			At:       logicgate.At,
		}
		if err = gate.CreateLogicGateNode(ctx, driver.Driver, dbname); err != nil {
			log.Fatalln(err)
		}
	}
	for at, relation := range parseresult.Nodes {
		// i1 <- lg, i2 <- lg, lg <- out
		// i1 <- lg
		if err := verpyverilog.LGtoIN(ctx, driver.Driver, dbname, relation.In1, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}

		// i2 <- lg
		if err := verpyverilog.LGtoIN(ctx, driver.Driver, dbname, relation.In2, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}

		// lg <- out
		if err := verpyverilog.OUTtoLG(ctx, driver.Driver, dbname, relation.Out, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}
	}
}

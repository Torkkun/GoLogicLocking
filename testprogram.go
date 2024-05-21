package main

import (
	"context"
	"goll/graph"
	"goll/parser"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

//テスト用の実行プログラム

func fullAddrtestMain() {
	testfilename := "test.txt"
	parseresult := parser.NewParse(testfilename)
	dbname := "neo4j"
	driver := graph.NewDriver()
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
		if err := LGtoIN(ctx, driver.Driver, dbname, relation.In1, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}

		// i2 <- lg
		if err := LGtoIN(ctx, driver.Driver, dbname, relation.In2, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}

		// lg <- out
		if err := OUTtoLG(ctx, driver.Driver, dbname, relation.Out, at, *parseresult.Declarations, parseresult.LogicGates); err != nil {
			log.Fatalln(err)
		}
	}
}

func isWire(wiremap map[string]parser.Wire, key string) (parser.Wire, bool) {
	val, iswire := wiremap[key]
	return val, iswire
}

func isIO(iomap map[string]parser.IOPort, key string) (parser.IOPort, bool) {
	val, isio := iomap[key]
	return val, isio
}

func LGtoIN(ctx context.Context, driver neo4j.DriverWithContext, dbname string, portname string, atnum int, declmaps parser.Decl, lgmap map[int]parser.LogicGate) error {
	ioval, bool := isIO(declmaps.IOPorts, portname)
	if bool {
		gio := graph.GateIO{
			Gate: graph.LogicGateNode{
				GateType: lgmap[atnum].GateType,
				At:       lgmap[atnum].At,
			},
			Io: graph.IONode{
				Type: string(ioval.Type),
				Name: ioval.Name,
			},
			At: atnum,
		}
		if err := gio.GatetoIO(ctx, driver, dbname); err != nil {
			return err
		}
		return nil
	}

	wireval, bool := isWire(declmaps.Wires, portname)
	if bool {
		gwire := graph.GateWire{
			Gate: graph.LogicGateNode{
				GateType: lgmap[atnum].GateType,
				At:       lgmap[atnum].At,
			},
			Wire: graph.WireNode{
				Name: wireval.Name,
			},
			At: atnum,
		}
		if err := gwire.GatetoWire(ctx, driver, dbname); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func OUTtoLG(ctx context.Context, driver neo4j.DriverWithContext, dbname string, portname string, atnum int, declmaps parser.Decl, lgmap map[int]parser.LogicGate) error {
	ioval, bool := isIO(declmaps.IOPorts, portname)
	if bool {
		gio := graph.GateIO{
			Gate: graph.LogicGateNode{
				GateType: lgmap[atnum].GateType,
				At:       lgmap[atnum].At,
			},
			Io: graph.IONode{
				Type: string(ioval.Type),
				Name: ioval.Name,
			},
			At: atnum,
		}
		if err := gio.IOtoGate(ctx, driver, dbname); err != nil {
			return err
		}
		return nil
	}

	wireval, bool := isWire(declmaps.Wires, portname)
	if bool {
		gwire := graph.GateWire{
			Gate: graph.LogicGateNode{
				GateType: lgmap[atnum].GateType,
				At:       lgmap[atnum].At,
			},
			Wire: graph.WireNode{
				Name: wireval.Name,
			},
			At: atnum,
		}
		if err := gwire.WiretoGate(ctx, driver, dbname); err != nil {
			return err
		}
		return nil
	}
	return nil
}

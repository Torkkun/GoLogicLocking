package cmd

import (
	"context"
	"fmt"
	"goll/graph"
	"goll/graph/verpyverilog"
	"goll/logiclocking"
	"goll/parser"
	"goll/parser/yosysjson"
	"goll/utils"
	"goll/vgenerator"
	"goll/vgenerator/funcmap"
	"log"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
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
	key, err := logiclocking.XorLock(ctx, driver, "neo4j", 2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(key)
}

func TestGrphToVerilog(t *testing.T) {
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
	modulename := "test"
	data := ConvertToVGeneratorData(ctx, driver, "neo4j")
	fmt.Println(data)
	data.ModuleName = modulename

	//if err := vgenerator.NewGenerator(*data, false); err != nil {
	//	log.Fatalln(err)
	//}
}

// 1bittest ver verilog
func ConvertToVGeneratorData(ctx context.Context, driver neo4j.DriverWithContext, dbname string) *vgenerator.GenerateData {
	all, err := verpyverilog.GetAllNodes(ctx, driver, "neo4j")
	if err != nil {
		log.Fatalln(err)
	}
	newData := new(vgenerator.GenerateData)
	var portlist []string
	var portDecls []*funcmap.PortDecl
	for _, in := range all.Ios.In {
		portlist = append(portlist, in.ION.Name)
		portDecls = append(portDecls, &funcmap.PortDecl{
			PortType:   funcmap.Input,
			SignalName: in.ION.Name,
		})
	}
	for _, out := range all.Ios.Out {
		portlist = append(portlist, out.ION.Name)
		portDecls = append(portDecls, &funcmap.PortDecl{
			PortType:   funcmap.Output,
			SignalName: out.ION.Name,
		})
	}
	var netDecls []*funcmap.NetDecl
	for _, wire := range all.Ws {
		netDecls = append(netDecls, &funcmap.NetDecl{
			NetType: funcmap.Wire,
			NetName: wire.WN.Name,
		})
	}
	var gateDecls []*funcmap.AssignDecl

	for _, gate := range all.Lgs {
		pre, err := verpyverilog.GetAllPredecessors(ctx, driver, dbname, gate.ElementId)
		if err != nil {
			log.Fatalln(err)
		}
		suc, err := verpyverilog.GetAllSuccessorNodes(ctx, driver, dbname, gate.ElementId)
		if err != nil {
			log.Fatalln(err)
		}
		yosystype := Selector(gate.LGN.GateType)
		var connection interface{}
		if yosystype == yosysjson.NOT || yosystype == yosysjson.BUF {
			connection = funcmap.BuforNot{
				// 接続する際に直接ゲートロジックからゲートロジックへ接続しているので間にワイヤーを作成しなくてはならない
				Y: pre[0].NodeName,
				A: suc[0].NodeName,
			}
		} else {
			connection = funcmap.Logic{
				Y: pre[0].NodeName,
				A: suc[0].NodeName,
				B: suc[0].NodeName,
			}
		}
		gateDecls = append(gateDecls, &funcmap.AssignDecl{
			ExpressionType: yosystype,
			Connection:     connection,
		})
	}
	lockg, err := verpyverilog.GetAllLockGateNodes(ctx, driver, dbname)
	if err != nil {
		log.Fatalln(err)
	}
	for _, lg := range lockg {
		pre, err := verpyverilog.GetAllPredecessors(ctx, driver, dbname, lg.LockGateNode.ElementId)
		if err != nil {
			log.Fatalln(err)
		}
		suc, err := verpyverilog.GetAllSuccessorNodes(ctx, driver, dbname, lg.LockGateNode.ElementId)
		if err != nil {
			log.Fatalln(err)
		}
		yosystype := Selector(lg.LockGateNode.GateType)
		connection := funcmap.Logic{
			Y: pre[0].NodeName,
			A: suc[0].NodeName,
			B: suc[0].NodeName,
		}
		gateDecls = append(gateDecls, &funcmap.AssignDecl{
			ExpressionType: yosystype,
			Connection:     connection,
		})

	}

	newData.PortList = portlist
	newData.PortDecl = portDecls
	newData.NetDecl = netDecls

	return newData
}

// yosys形式のタイプにするために一時的な置き換える関数
func Selector(logictype string) string {
	switch logictype {
	case string(utils.And):
		return yosysjson.AND
	case string(utils.Nand):
		return yosysjson.NAND
	case string(utils.Or):
		return yosysjson.OR
	case string(utils.Nor):
		return yosysjson.NOR
	case string(utils.Xor):
		return yosysjson.XOR
	case string(utils.Xnor):
		return yosysjson.XNOR
	case string(utils.Not):
		return yosysjson.NOT
	default:
		log.Fatalln("This Type is not implement")
		return ""
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

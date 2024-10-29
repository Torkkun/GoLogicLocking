package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goll/circuit"
	"goll/graph"
	"goll/graph/veryosys"
	"goll/logiclocking/veryosysll"
	"goll/parser/yosysjson"
	"log"
	"os"
	"testing"
)

// key is num, value is {IO or NET or CELL} and element id
type Mapv struct {
	Type      string
	ElementId string
}

// 個別テスト用のもの
type testYosys struct {
	module *yosysjson.Module
	driver *graph.GraphDB
}

func SetUp() *testYosys {
	var path = "C:\\Users\\onigi\\projects\\GoLogicLocking\\testverilog\\workspace\\fulladd\\testfulladd.json"
	var topmodule = "fulladd"
	driver := graph.NewDriver()

	file, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	jsondata := new(yosysjson.YosysJson)
	if err = json.Unmarshal(file, jsondata); err != nil {
		log.Fatalln(err)
	}
	fulladdmod := jsondata.Modules[topmodule]
	return &testYosys{
		module: &fulladdmod,
		driver: driver,
	}
}

func TestPort(t *testing.T) {
	conf := SetUp()
	ctx := context.Background()

	defer conf.driver.Driver.Close(ctx)
	tmpmap := make(map[int]*Mapv)
	_, err := PortCreateTest(t, conf, ctx, tmpmap)
	if err != nil {
		TestDeleteAllDB(t)
		log.Fatalln("Port to Graph Error")
	}
}

func PortCreateTest(t *testing.T, conf *testYosys, ctx context.Context, tmpmap map[int]*Mapv) (map[int]*Mapv, error) {
	gport := yosysjson.PortsConvertToStoreGraph(conf.module.Ports)
	for _, v := range gport {
		fmt.Println(v)
	}

	for _, io := range gport {
		elementid, err := veryosys.CreateInOutNode(ctx, conf.driver.Driver, conf.driver.DBname, io)

		if err != nil {
			TestDeleteAllDB(t)
			log.Fatalf("%v", err.Error())
		}

		tmpmap[io.BitNum] = &Mapv{
			Type:      "IO",
			ElementId: elementid,
		}
		fmt.Println(tmpmap)
	}
	if len(tmpmap) != len(gport) {
		err := errors.New("Port to Graph Error")
		return nil, err
	}
	return tmpmap, nil
}

// Net Wire
func TestNet(t *testing.T) {
	conf := SetUp()
	ctx := context.Background()

	defer conf.driver.Driver.Close(ctx)

	tmpmap := make(map[int]*Mapv)

	_, err := NetCreateTest(t, conf, ctx, tmpmap)
	if err != nil {
		TestDeleteAllDB(t)
		log.Fatalln(err.Error())
	}
}

func NetCreateTest(t *testing.T, conf *testYosys, ctx context.Context, tmpmap map[int]*Mapv) (map[int]*Mapv, error) {
	gname := yosysjson.NetsConvertToStoreGraph(conf.module.NetName)
	for _, v := range gname {
		fmt.Println(v)
	}

	// テスト用のノードカウント変数
	var nodecount int
	// テスト用の差分の値
	diffnum := len(tmpmap)

	for _, net := range gname {
		newnet := new(veryosys.DBNetName)
		newnet.Netname = net.Netname
		newnet.Attributes.Src = net.Attributes.Src

		nodecount += len(net.Bits)

		for _, bitnum := range net.Bits {
			newnet.BitNum = bitnum
			elementid, err := veryosys.CreateWireNode(ctx, conf.driver.Driver, conf.driver.DBname, newnet)

			if err != nil {
				TestDeleteAllDB(t)
				log.Fatalln(err.Error())
			}

			tmpmap[bitnum] = &Mapv{
				Type:      "Wire",
				ElementId: elementid,
			}
			fmt.Println(tmpmap)
		}
	}
	if len(tmpmap)-diffnum != nodecount {
		err := errors.New("Net to Graph Error:")
		return nil, err
	}
	log.Println("Create Wire Node OK")
	return tmpmap, nil
}

func TestCell(t *testing.T) {
	conf := SetUp()
	ctx := context.Background()

	defer conf.driver.Driver.Close(ctx)

	_, err := CellCreateTest(t, conf, ctx)
	if err != nil {
		log.Fatalln(err.Error())
		TestDeleteAllDB(t)
	}
}

type celltest struct {
	Type string //Cell
}

func CellCreateTest(t *testing.T, conf *testYosys, ctx context.Context) (map[string]*celltest, error) {
	gcell, gconn := yosysjson.CellsConvertToStoreGraph(conf.module.Cells)
	for _, v := range gcell {
		fmt.Println(v)
	}

	tmpmap := make(map[string]*celltest) // MapVとは別

	for _, cell := range gcell {
		elementid, err := veryosys.CreateCellNode(ctx, conf.driver.Driver, conf.driver.DBname, cell)
		if err != nil {
			TestDeleteAllDB(t)
			log.Fatalln(err.Error())
		}

		tmpmap[elementid] = &celltest{
			Type: elementid,
		}
		fmt.Println(tmpmap)
	}

	// コネクション
	for _, conn := range gconn {
		for _, v := range conn {
			fmt.Println(v)
		}
	}

	if len(tmpmap) != len(gcell) {
		err := errors.New("Cell to Graph Error")
		return nil, err
	}
	log.Println("Create Cell Node OK")
	return tmpmap, nil
}

// DBへ全ノードの作成を行う
func TestConnection(t *testing.T) {
	conf := SetUp()
	ctx := context.Background()

	defer conf.driver.Driver.Close(ctx)

	tmpmap := make(map[int]*Mapv)

	// テストのノード作成を行う
	tmpmap, err := PortCreateTest(t, conf, ctx, tmpmap)
	if err != nil {
		log.Fatalln(err.Error())
	}
	tmpmap, err = NetCreateTest(t, conf, ctx, tmpmap)
	if err != nil {
		log.Fatalln(err.Error())
	}
	gcell, gconns := yosysjson.CellsConvertToStoreGraph(conf.module.Cells)
	for _, v := range gcell {
		fmt.Println(v)
	}
	log.Println("セルのコンバート完了")
	// 推定されるコネクションの数
	var connNum int
	// 作成した数
	var cnum int

	cellnodes := make(map[string]string) // K AttrSrc, V elementId
	for _, cell := range gcell {
		elementid, err := veryosys.CreateCellNode(ctx, conf.driver.Driver, conf.driver.DBname, cell)
		if err != nil {
			log.Fatalf("Create Cell Node Error: %v", err.Error())
		}
		cellnodes[cell.Attributes.Src] = elementid
	}

	for src, conn := range gconns {
		connNum += len(conn) // connectioの数

		elementid := cellnodes[src]
		for _, v := range conn {
			p := tmpmap[v.BitNum]
			if v.Type == "input" {
				err := veryosys.CellConnection(ctx, conf.driver.Driver, conf.driver.DBname, &veryosys.ConnectionPair{
					Predecessor: veryosys.Node{
						Type:      "Cell",
						ElementId: elementid,
					},
					Successor: veryosys.Node{
						Type:      p.Type,
						ElementId: p.ElementId,
					},
				})
				if err != nil {
					log.Fatalf("Create Connection Node Error: %v", err.Error())
				}
				// 作成完了につき　＋１
				log.Printf("CellID:%v -> %v", elementid, v.BitNum)
				cnum += 1
			} else if v.Type == "output" {
				// Netで作成したものとCellを接続する
				err := veryosys.CellConnection(ctx, conf.driver.Driver, conf.driver.DBname, &veryosys.ConnectionPair{
					Predecessor: veryosys.Node{
						Type:      p.Type,
						ElementId: p.ElementId,
					},
					Successor: veryosys.Node{
						Type:      "Cell",
						ElementId: elementid,
					},
				})
				if err != nil {
					log.Fatalf("Create Connection Node Error: %v", err.Error())
				}
				// 作成完了につき　＋１
				log.Printf("CellID:%v -> %v", v.BitNum, elementid)
				cnum += 1
			}
		}
	}
	// 作成した数を確かめる
	if connNum == cnum {
		log.Printf("Connection Created OK: created %v", connNum)

	} else {
		log.Fatalln("Connection Create Faled")
	}
}

func TestCreateGraph(t *testing.T) {
	conf := SetUp()
	ctx := context.Background()

	defer conf.driver.Driver.Close(ctx)
	if err := circuit.TranslationCircuit(ctx, conf.driver.Driver, conf.driver.DBname, conf.module); err != nil {
		log.Fatalln(err)
	}
}

func TestRandomXorXnorWrpTx(t *testing.T) {
	conf := SetUp()
	ctx := context.Background()

	defer conf.driver.Driver.Close(ctx)

	keys, err := veryosysll.RandomXorXnorWrpTxTest(ctx, conf.driver.Driver, conf.driver.DBname, 2)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(keys)
}

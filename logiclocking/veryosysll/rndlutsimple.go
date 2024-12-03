package veryosysll

import (
	"context"
	"fmt"
	"goll/graph/veryosys"
	"goll/utils"
	"log"
	"math"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
	"golang.org/x/exp/maps"
)

var testnodes = []dbtype.Node{
	{
		Id:        4,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:4",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:6.16-6.21",
			"type":    "$_AND_",
		},
	},
	{
		Id:        5,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:5",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:6.26-6.33",
			"type":    "$_AND_",
		},
	},
	{
		Id:        6,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:6",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:6.15-6.34",
			"type":    "$_OR_",
		},
	},
	{
		Id:        7,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:7",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:6.38-6.45",
			"type":    "$_AND_",
		},
	},
	{
		Id:        8,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:8",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:6.15-6.46",
			"type":    "$_OR_",
		},
	},
	{
		Id:        9,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:9",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:5.12-5.17",
			"type":    "$_XOR_",
		},
	},
	{
		Id:        10,
		ElementId: "4:af7b1260-e76d-49b6-aeb5-0f18dc98a01e:10",
		Labels:    []string{"Cell"},
		Props: map[string]interface{}{
			"attrsrc": "fulladd.v:5.12-5.23",
			"type":    "$_XOR_",
		},
	},
}

// Assumptionというすべての組み合わせから特定のノードに対して設定される仮定の真偽値
func TestAssumptionExec(ctx context.Context, driver neo4j.DriverWithContext, dbname string, numgates, lutwidth int) {
	types := []string{veryosys.IONode, veryosys.WireNode}
	// potential gates
	nodes, err := veryosys.GetNodes(ctx, driver, dbname, types)
	if err != nil {
		log.Fatalln(err)
	}
	if len(nodes) < numgates {
		err = fmt.Errorf("ERROR:to much keylength, plese input less than %v", len(nodes))
		log.Fatalln(err)
	}
	potential := convertToPotentialList(nodes)

	gates := utils.RandomNumbers(len(nodes), numgates)

	var tmplist []string                // ランダムにnumgate分選択されたロックするゲートの個数
	tmpmap := make(map[string][]string) // ロックするゲートの子mapにすることで自動的にSetする
	for _, gate := range gates {
		node := nodes[gate]
		elementid := node.GetElementId()
		tmplist = append(tmplist, elementid)
		pre, err := veryosys.GetPredecessorNodes(ctx, driver, dbname, elementid)
		if err != nil {
			log.Fatalln(err)
		}
		// transitive fanout: ゲートにつながる親を取得
		for _, node := range pre.Nodes {
			tmpmap[node.GetElementId()] = append(tmpmap[node.GetElementId()], elementid)
		}
	}

	// nodeから選択したゲートを削除
	potential = rejectNodes(potential, tmplist)

	// nodeからgatesの子のノードを削除
	successors := maps.Keys(tmpmap)
	potential = rejectNodes(potential, successors)

	// lut用のMUXの設定を用意
	//muxwidth := 2 * lutwidth

	// testkeylist
	var keylist []string

	for i, indxnum := range gates {
		nodeID := nodes[indxnum].GetElementId()
		// fanin
		suc, err := veryosys.GetSuccessorNodes(ctx, driver, dbname, nodeID)
		if err != nil {
			log.Fatalln(err)
		}
		var padding []string
		var faninlist []string
		isPadding := len(suc.Nodes) <= lutwidth
		// lutwidth-len(suc.Nodes) == 0の場合でもやる
		if isPadding {
			for _, node := range suc.Nodes {
				faninlist = append(faninlist, node.GetElementId())
			}
			tmpmap := rejectNodes(potential, faninlist)
			//fmt.Println(tmpmap)
			tmppotential := maps.Keys(tmpmap)

			fmt.Println(suc.Nodes)
			fmt.Println("padding:", lutwidth-len(suc.Nodes))

			paddingidx := utils.RandomNumbers(len(tmppotential), lutwidth-len(suc.Nodes))

			fmt.Println(paddingidx)

			for _, pidx := range paddingidx {
				padding = append(padding, tmppotential[pidx])
			}
		} else {
			log.Panicln("could not find enough viable gates for padding")
		}
		fmt.Println("padding list(ないときはない):", padding)
		// connect keys
		numvars := len(faninlist) + len(padding)
		totalCombinations := int(math.Pow(2, float64(numvars)))

		// 組み合わせの作成
		for j := 0; j < totalCombinations; j++ {
			vs := make([]bool, numvars)
			for k := 0; k < numvars; k++ {
				vs[k] = (j>>k)&1 == 1
			}
			fmt.Println(vs)

			// assumptionsの作成
			assumptions := make(map[string]bool)
			for idx, signal := range append(faninlist, padding...) {
				if contains(faninlist, signal) {
					assumptions[signal] = vs[numvars-1-idx]
				}
			}
			fmt.Println(assumptions)

			// 動的キー生成
			keyIn := fmt.Sprintf("key_%d", i*int(math.Pow(2, float64(lutwidth)))+j)
			keylist = append(keylist, keyIn)
		}
	}
	fmt.Println(keylist)

}

// Assumptionを作成した後にその仮定に対してSATを行い制約として元の回路を与え、仮定（assumption）に基づき解があるかどうか判定する
func TestSATSolverExec(ctx context.Context, driver neo4j.DriverWithContext, dbname string, numgates, lutwidth int) {
	types := []string{veryosys.IONode, veryosys.WireNode}
	// potential gates
	nodes, err := veryosys.GetNodes(ctx, driver, dbname, types)
	if err != nil {
		log.Fatalln(err)
	}
	if len(nodes) < numgates {
		err = fmt.Errorf("ERROR:to much keylength, plese input less than %v", len(nodes))
		log.Fatalln(err)
	}
	potential := convertToPotentialList(nodes)

	gates := utils.RandomNumbers(len(nodes), numgates)

	var tmplist []string                // ランダムにnumgate分選択されたロックするゲートの個数
	tmpmap := make(map[string][]string) // ロックするゲートの子mapにすることで自動的にSetする
	for _, gate := range gates {
		node := nodes[gate]
		elementid := node.GetElementId()
		tmplist = append(tmplist, elementid)
		pre, err := veryosys.GetPredecessorNodes(ctx, driver, dbname, elementid)
		if err != nil {
			log.Fatalln(err)
		}
		// transitive fanout: ゲートにつながる親を取得
		for _, node := range pre.Nodes {
			tmpmap[node.GetElementId()] = append(tmpmap[node.GetElementId()], elementid)
		}
	}

	// nodeから選択したゲートを削除
	potential = rejectNodes(potential, tmplist)

	// nodeからgatesの子のノードを削除
	successors := maps.Keys(tmpmap)
	potential = rejectNodes(potential, successors)

	// lut用のMUXの設定を用意
	//muxwidth := 2 * lutwidth

	// testkeylist
	var keylist []string

	for i, indxnum := range gates {
		nodeID := nodes[indxnum].GetElementId()
		// fanin
		suc, err := veryosys.GetSuccessorNodes(ctx, driver, dbname, nodeID)
		if err != nil {
			log.Fatalln(err)
		}
		var padding []string
		var faninlist []string
		isPadding := len(suc.Nodes) <= lutwidth
		// lutwidth-len(suc.Nodes) == 0の場合でもやる
		if isPadding {
			for _, node := range suc.Nodes {
				faninlist = append(faninlist, node.GetElementId())
			}
			tmpmap := rejectNodes(potential, faninlist)
			//fmt.Println(tmpmap)
			tmppotential := maps.Keys(tmpmap)

			fmt.Println(suc.Nodes)
			fmt.Println("padding:", lutwidth-len(suc.Nodes))

			paddingidx := utils.RandomNumbers(len(tmppotential), lutwidth-len(suc.Nodes))

			fmt.Println(paddingidx)

			for _, pidx := range paddingidx {
				padding = append(padding, tmppotential[pidx])
			}
		} else {
			log.Panicln("could not find enough viable gates for padding")
		}
		fmt.Println("padding list(ないときはない):", padding)
		// connect keys
		numvars := len(faninlist) + len(padding)
		totalCombinations := int(math.Pow(2, float64(numvars)))

		// 組み合わせの作成
		for j := 0; j < totalCombinations; j++ {
			vs := make([]bool, numvars)
			for k := 0; k < numvars; k++ {
				vs[k] = (j>>k)&1 == 1
			}
			fmt.Println(vs)

			// assumptionsの作成
			assumptions := make(map[string]bool)
			for idx, signal := range append(faninlist, padding...) {
				if contains(faninlist, signal) {
					assumptions[signal] = vs[numvars-1-idx]
				}
			}
			fmt.Println(assumptions)

			// 動的キー生成
			keyIn := fmt.Sprintf("key_%d", i*int(math.Pow(2, float64(lutwidth)))+j)
			keylist = append(keylist, keyIn) // testなので後で消す

			// SAT
			//実行時エラーはエラーが返る
			//errではなくnilの場合はkey[key_in] = False
			//それ以外は result[gate]になるようにする

		}
	}
	fmt.Println(keylist)

}

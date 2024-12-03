package veryosysll

import (
	"context"
	"fmt"
	"goll/graph/veryosys"
	"goll/utils"
	"math"
	"slices"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"golang.org/x/exp/maps"
)

// randomLogicLockingのロックゲート部分のみLUT化させる
func TMPRandomLUT(ctx context.Context, driver neo4j.DriverWithContext, dbname string, numgates, lutwidth int) (map[string]bool, error) {
	// num gates: ロックするゲートの数
	// lut width: LUTの幅ロックされたゲートの最大のfaninを定義する

	// ペアではなくてゲートを取得するものに変更する
	types := []string{veryosys.IONode, veryosys.WireNode}
	// potential gates
	nodes, err := veryosys.GetNodes(ctx, driver, dbname, types)
	if err != nil {
		return nil, err
	}
	if len(nodes) < numgates {
		err = fmt.Errorf("ERROR:to much keylength, plese input less than %v", len(nodes))
		return nil, err
	}
	potential := convertToPotentialList(nodes)

	// ランダムにロックするゲートを選ぶ
	// リレーションの数を範囲としてランダムにシャッフルしkeylen分だけ先頭からスライスする
	gates := utils.RandomNumbers(len(nodes), numgates)

	var tmplist []string                // ランダムにnumgate分選択されたロックするゲートの個数
	tmpmap := make(map[string][]string) // ロックするゲートの子mapにすることで自動的にSetする
	for _, gate := range gates {
		node := nodes[gate]
		elementid := node.GetElementId()
		tmplist = append(tmplist, elementid)
		pre, err := veryosys.GetPredecessorNodes(ctx, driver, dbname, elementid)
		if err != nil {
			return nil, err
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
	muxwidth := 2 * lutwidth

	// TODO: KEY USED
	//key := make(map[string]bool)

	// session 作成
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: dbname})
	defer session.Close(ctx)
	// 明示的なTransactionを開始
	_, err = veryosys.Transaction(ctx, session, dbname, func(tx neo4j.ExplicitTransaction) (interface{}, error) {
		// keyサイズとLUTサイズを選択された数分それぞれ作成する
		for i, indxnum := range gates {
			nodeID := nodes[indxnum].GetElementId()
			// fanout
			//pre, err := veryosys.GetPredecessorNodes(ctx, driver, dbname, nodeID)
			//if err != nil {
			//	return nil, err
			//}

			// fanin
			suc, err := veryosys.GetSuccessorNodes(ctx, driver, dbname, nodeID)
			if err != nil {
				return nil, err
			}

			// TODO padding
			// ゲートのfaninが lut width 以下の場合ダミーの接続を追加する
			// ここはよくわからん
			var padding []string
			var faninlist []string
			isPadding := len(suc.Nodes) <= lutwidth
			if isPadding {
				for _, node := range suc.Nodes {
					faninlist = append(faninlist, node.GetElementId())
				}
				tmpmap := rejectNodes(potential, faninlist)
				tmppotential := maps.Keys(tmpmap)
				padding, err = utils.Sample(tmppotential, lutwidth-len(suc.Nodes))
				if err != nil {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("could not find enough viable gates for padding")
			}

			// create LUT
			lutid, err := AddLUTNode(tx, ctx, fmt.Sprintf("lut_%d", i), muxwidth)
			if err != nil {
				return nil, err
			}

			// connect keys
			product := utils.ProductBool(len(suc.Nodes) - len(faninlist))
			// 組み合わせの作成
			var keylists [][]bool
			for {
				keylist := product()
				if len(keylist) == 0 {
					break
				}
				keylists = append(keylists, keylist)
			}
			// assumptionsの作成
			tmpfaninandpadding := append(faninlist, padding...)

			for j, vs := range keylists {
				assumptions := make(map[string]bool)
				for idx, v := range slices.Backward(vs) {
					if contains(faninlist, tmpfaninandpadding[idx]) {
						assumptions[tmpfaninandpadding[idx]] = v
					}
				}

				// 動的キー生成
				keyIn := fmt.Sprintf("key_%d", i*int(math.Pow(2, float64(lutwidth)))+j)

				// いったん概念的なグラフ上のLUTに接続する?
				newkeyid, err := veryosys.CreateInOutNodeTx(tx, ctx, &veryosys.Port{
					Direction: "input",
					Name:      keyIn,
				})
				if err != nil {
					return nil, err
				}
				if err := veryosys.CellConnectionTx(tx, ctx, &veryosys.ConnectionPair{
					Predecessor: veryosys.Node{
						Type:      "LLLut",
						ElementId: lutid,
					},
					Successor: veryosys.Node{
						Type:      "IO",
						ElementId: newkeyid,
					},
				}); err != nil {
					return nil, err
				}

				// SATソルバ
				//TODO
				// 下は例の実装
				// result(bool) = satsolver()
				//if nil == result
				// key[keyIn] = False
				//else
				// key[keyIn] = result[gate]
			}
			// TODO: 接続の変更
			// 接続を削除
			// 接続をコネクト
			// 元のゲートを削除

		}

		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func convertToPotentialList(nodes []neo4j.Node) map[string]int {
	pgate := make(map[string]int)
	for i, node := range nodes {
		pgate[node.GetElementId()] = i
	}
	return pgate
}

// もととなるリストと消去するノードのID
func rejectNodes(potentialGates map[string]int, elementIds []string) map[string]int {
	for _, id := range elementIds {
		delete(potentialGates, id)
	}
	return potentialGates
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

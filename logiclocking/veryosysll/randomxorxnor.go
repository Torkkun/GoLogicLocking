package veryosysll

import (
	"context"
	"fmt"
	"goll/graph/veryosys"
	"goll/utils"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func RandomXorXnorWrpTxTest(ctx context.Context, driver neo4j.DriverWithContext, dbname string, keylen int) (map[string]bool, error) {
	pair, err := veryosys.GetRelationshipsAndPairNodes(ctx, driver, dbname)
	if err != nil {
		return nil, err
	}

	if len(pair.Relation) < keylen {
		err = fmt.Errorf("ERROR:to much keylength, plese input less than %v", len(pair.Relation))
		return nil, err
	}
	// リレーションの数を範囲としてランダムにシャッフルしkeylen分だけ先頭からスライスする
	rnum := utils.RandomNumbers(len(pair.Relation), keylen)
	key := make(map[string]bool)

	// session 作成
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: dbname})
	defer session.Close(ctx)
	// 明示的なTransactionを開始
	_, err = veryosys.Transaction(ctx, session, dbname, func(tx neo4j.ExplicitTransaction) (interface{}, error) {
		for i, indxnum := range rnum {
			b, err := utils.Choice([]bool{true, false})
			if err != nil {
				return nil, err
			}
			keystring := fmt.Sprintf("key_%v", i)
			key[keystring] = b
			var gateType string
			if b {
				gateType = "XNOR"
			} else {
				gateType = "XOR"
			}
			prenode := new(veryosys.Node)
			prenode.ElementId = pair.Pre[indxnum].GetElementId()
			prenode.Type = pair.Pre[indxnum].Labels[0]

			sucnode := new(veryosys.Node)
			sucnode.ElementId = pair.Suc[indxnum].GetElementId()
			sucnode.Type = pair.Suc[indxnum].Labels[0]

			relationId := pair.Relation[indxnum].GetElementId()

			// ロックゲートの追加
			newL := new(veryosys.LockGateNode)
			newL.GateType = gateType
			newL.LockType = "randomll"

			newelementid, err := veryosys.CreateRandomLLGateNodeTx(tx, ctx, newL)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			// ロックしたキーを設定
			newio := new(veryosys.Port)
			newio.Name = fmt.Sprintf("key_%v", i)
			newio.Direction = "input"
			newio.BitWidth = 1 // BitNumは無し
			newinleyid, err := veryosys.CreateInOutNodeTx(tx, ctx, newio)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			// もともとのロジックゲートのOUTを外して付け替える
			// Relation削除
			if err := veryosys.DeleteConnectionTx(tx, ctx, relationId); err != nil {
				log.Println(err)
				return nil, err
			}
			// xor/xnorと親ノードを接続
			if err := veryosys.CellConnectionTx(tx, ctx, &veryosys.ConnectionPair{
				Predecessor: veryosys.Node{
					Type:      prenode.Type,
					ElementId: prenode.ElementId,
				},
				Successor: veryosys.Node{
					Type:      "LLCell",
					ElementId: newelementid,
				},
			}); err != nil {
				log.Println(err)
				return nil, err
			}
			// xor/xnorと子ノードを接続
			if err := veryosys.CellConnectionTx(tx, ctx, &veryosys.ConnectionPair{
				Predecessor: veryosys.Node{
					Type:      "LLCell",
					ElementId: newelementid,
				},
				Successor: veryosys.Node{
					Type:      sucnode.Type,
					ElementId: sucnode.ElementId,
				},
			}); err != nil {
				log.Println(err)
				return nil, err
			}
			// 解除キーと接続
			if err := veryosys.CellConnectionTx(tx, ctx, &veryosys.ConnectionPair{
				Predecessor: veryosys.Node{
					Type:      "LLCell",
					ElementId: newelementid,
				},
				Successor: veryosys.Node{
					Type:      "IO",
					ElementId: newinleyid,
				},
			}); err != nil {
				log.Println(err)
				return nil, err
			}
		}
		// Transaction 終了
		return nil, nil
	})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return key, nil
}

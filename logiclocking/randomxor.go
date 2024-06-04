package logiclocking

import (
	"context"
	"fmt"
	"goll/graph"
	"goll/utils"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// TODO
func XorLock(ctx context.Context, driver neo4j.DriverWithContext, dbname string, keylen int) error {
	// db側でコピー（Cyperクエリのバックアップクエリ）

	// ノード数 - アウトプットの数でロック用のゲートを作成
	//  - （ランダムにXORとXNORのゲート情報を作成していると思われる）

	// 上のゲートをひとつづつ取り出し、キーリスト（map）にてTrueまたはFalseを選択
	//  - circuitgraphのchoicesとsample関数をチラ見する
	all, err := graph.GetAllNodes(ctx, driver, dbname)
	if err != nil {
		return err
	}
	// GateLogicとWireとInputからランダムにロックするゲートを選ぶ
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

	gates, err := utils.Sample(idlist, keylen)
	if err != nil {
		log.Fatalln(err)
	}

	for _, gate := range gates {
		// TRUEまたはFALSEがゲートに設定されるが、
		//  - TRUEならXNOR、FALSEならXORとゲートタイプが決定される
		fmt.Println(gate)
		// ゲートの追加とインプットの追加を行う
		//　もともとのロジックゲートのOUTと
	}

	return nil
}

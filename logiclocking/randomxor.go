package logiclocking

import (
	"context"
	"fmt"
	"goll/graph/verpyverilog"
	"goll/utils"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Id struct {
	ElementId string
	Id        int64
}

func XorLock(ctx context.Context, driver neo4j.DriverWithContext, dbname string, keylen int) (map[string]bool, error) {
	// db側でコピー（Cyperクエリのバックアップクエリ）

	all, err := verpyverilog.GetAllNodes(ctx, driver, dbname)
	if err != nil {
		return nil, err
	}
	// GateLogicとWireとInputからランダムにロックするゲートを選ぶ

	var idlist []Id
	for _, v := range all.Lgs {
		idlist = append(idlist, Id{
			ElementId: v.ElementId,
			Id:        v.Id,
		})
	}

	if len(idlist) < keylen {
		err = fmt.Errorf("ERROR:to much keylength, plese input less than %v", len(idlist))
		return nil, err
	}

	gates, err := utils.Sample(idlist, keylen)
	if err != nil {
		return nil, err
	}
	key := make(map[string]bool)
	for i, gate := range gates {
		// TRUEまたはFALSEがゲートに設定される
		b, err := utils.Choice([]bool{true, false})
		if err != nil {
			return nil, err
		}
		keystring := fmt.Sprintf("key_%v", i)
		key[keystring] = b
		//  - TRUEならXNOR、FALSEならXORとゲートタイプが決定される
		var gateType utils.GateType
		if b {
			gateType = utils.Xnor
		} else {
			gateType = utils.Xor
		}

		pres, err := GetGatePredecessorNodes(ctx, driver, dbname, gate)
		if err != nil {
			return nil, err
		}
		// コネクションを外す
		// 1 対 1で繋がっているものとする
		// ゲートの追加とインプットの追加を行う
		newL := new(verpyverilog.LockGateNode)
		newL.GateType = string(gateType)
		newL.LockType = "randomll"
		newL.Name = fmt.Sprintf("key_gate_%v", i)
		newLGEid, err := newL.CreateLockingGateNode(ctx, driver, dbname)
		if err != nil {
			return nil, err
		}
		// ノード作成後ElementId取得して設定
		newL.ElementId = newLGEid

		newio := new(verpyverilog.IONode)
		newio.Name = fmt.Sprintf("key_%v", i)
		newio.Type = "IN"
		if err = newio.CreateInOutNode(ctx, driver, dbname); err != nil {
			return nil, err
		}
		// もともとのロジックゲートのOUTを外して、付け替える
		if len(pres.IOaR) != 0 {
			ioar := pres.IOaR[0]
			if err = verpyverilog.DeleteRelationIOtoGateByElementId(ctx, driver, dbname, ioar.Relation.ElementId); err != nil {
				return nil, err
			}
			// xorとio(out)を接続
			funout := ioar.Neo4JIO
			newIOtoLG := new(verpyverilog.LockGateIO)
			newIOtoLG.Gate = newL
			newIOtoLG.Io = funout.ION
			if err = newIOtoLG.IOtoLLGateByElementId(ctx, driver, dbname); err != nil {
				return nil, err
			}

		} else if len(pres.WaR) != 0 {
			war := pres.WaR[0]
			if err = verpyverilog.DeleteRelationWiretoGateByElementId(ctx, driver, dbname, war.Relation.ElementId); err != nil {
				return nil, err
			}
			// xorとwireを接続
			funout := war.Neo4JWire
			newWtoLG := new(verpyverilog.LockGateWire)
			newWtoLG.Gate = newL
			newWtoLG.Wire = funout.WN
			if err = newWtoLG.WiretoLLGateByElementId(ctx, driver, dbname); err != nil {
				return nil, err
			}

		} else {
			err = fmt.Errorf("GateRelation is missing")
			return nil, err
		}
		// key追加とxorゲートと接続
		newLGtoIO := new(verpyverilog.LockGateIO)
		newLGtoIO.Gate = newL
		newLGtoIO.Io = newio
		if err = newLGtoIO.LLGatetoIOByElementId(ctx, driver, dbname); err != nil {
			return nil, err
		}
		// xorゲートと元のゲートと接続
		funin := all.Lgs[gate.ElementId].LGN
		newLgtoG := new(verpyverilog.LLGateGate)
		newLgtoG.LGN = funin
		newLgtoG.LLGN = newL
		if err = newLgtoG.LLGatetoGateElementId(ctx, driver, dbname); err != nil {
			return nil, err
		}
	}
	return key, nil
}

type GatePredecessorNodes struct {
	WaR  []*verpyverilog.GetNeo4JWireAndRelation
	IOaR []*verpyverilog.GetNeo4JIoAndRelation
}

func GetGatePredecessorNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string, idlist Id) (*GatePredecessorNodes, error) {
	war, err := verpyverilog.GetWiretoGateRelationByElementId(ctx, driver, dbname, idlist.ElementId)
	if err != nil {
		return nil, err
	}
	ioar, err := verpyverilog.GetIOtoGateRelationByElementId(ctx, driver, dbname, idlist.ElementId)
	if err != nil {
		return nil, err
	}
	return &GatePredecessorNodes{
		WaR:  war,
		IOaR: ioar,
	}, nil
}

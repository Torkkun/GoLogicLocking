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
		newLGEid, err := newL.CreateLockingGateNode(ctx, driver, dbname)
		if err != nil {
			return nil, err
		}

		// ロックしたキーを設定
		newio := new(verpyverilog.IONode)
		newio.Name = fmt.Sprintf("key_%v", i)
		newio.Type = "IN"
		newinkeyid, err := newio.CreateInOutNode(ctx, driver, dbname)
		if err != nil {
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
			newIOtoLG.GateElementId = newLGEid
			newIOtoLG.IoElementId = funout.ElementId
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
			newWtoLG.GateElementId = newLGEid
			newWtoLG.WireElementId = funout.ElementId
			if err = newWtoLG.WiretoLLGateByElementId(ctx, driver, dbname); err != nil {
				return nil, err
			}
		} else {
			err = fmt.Errorf("GateRelation is missing")
			return nil, err
		}
		// key追加とxorゲートと接続
		newLGtoIO := new(verpyverilog.LockGateIO)
		newLGtoIO.GateElementId = newLGEid
		newLGtoIO.IoElementId = newinkeyid
		if err = newLGtoIO.LLGatetoIOByElementId(ctx, driver, dbname); err != nil {
			return nil, err
		}
		// xorゲートと元のゲートと接続
		// その際に間にwireを介する
		newconnwire := new(verpyverilog.WireNode)
		newconnwire.Name = fmt.Sprintf("key_gate_%d", i)
		newwireid, err := newconnwire.CreateWireNode(ctx, driver, dbname)
		if err != nil {
			return nil, err
		}
		// 元のLogicGate <- 新しく作成したwire
		gw := new(verpyverilog.GateWire)
		gw.GateElementId = gate.ElementId
		gw.WireElementId = newwireid
		if err = gw.WiretoGateByElementId(ctx, driver, dbname); err != nil {
			return nil, err
		}
		// 新しく作成したwire <- LockLogicGate
		llgw := new(verpyverilog.LockGateWire)
		llgw.GateElementId = newLGEid
		llgw.WireElementId = newwireid
		if err = llgw.LLGatetoWireByElementId(ctx, driver, dbname); err != nil {
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

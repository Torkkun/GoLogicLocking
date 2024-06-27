package verpyverilog

import (
	"context"
	"fmt"
	"goll/parser"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// ワイヤーかどうか存在するかどうか
// 廃止予定
func IsWire(wiremap map[string]parser.Wire, key string) (parser.Wire, bool) {
	val, iswire := wiremap[key]
	return val, iswire
}

// 廃止予定
// 存在するかどうか
func IsIO(iomap map[string]parser.IOPort, key string) (parser.IOPort, bool) {
	val, isio := iomap[key]
	return val, isio
}

/* func LGtoIN(ctx context.Context, driver neo4j.DriverWithContext, dbname string, gateElementId string) error {

	return nil
}

func OUTtoLG(ctx context.Context, driver neo4j.DriverWithContext, dbname string, portname string, atnum int, declmaps parser.Decl, lgmap map[int]parser.LogicGate) error {
	ioval, bool := IsIO(declmaps.IOPorts, portname)
	if bool {
		gio := GateIO{
			Gate: &LogicGateNode{
				GateType: string(lgmap[atnum].GateType),
			},
			Io: &IONode{
				Type: string(ioval.Type),
				Name: ioval.Name,
			},
		}
		if err := gio.IOtoGate(ctx, driver, dbname); err != nil {
			return err
		}
		return nil
	}

	wireval, bool := IsWire(declmaps.Wires, portname)
	if bool {
		gwire := GateWire{
			Gate: &LogicGateNode{
				GateType: string(lgmap[atnum].GateType),
			},
			Wire: &WireNode{
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
} */

// nep4jのElementidをKeyとする
type AllNode struct {
	Ios *IONodes
	Lgs map[string]*GetNeo4JLogicNode
	Ws  map[string]*GetNeo4JWireNode
}

type IONodes struct {
	In  map[string]*GetNeo4JIONode
	Out map[string]*GetNeo4JIONode
}

func GetAllNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (*AllNode, error) {
	ionodes, err := GetAllIONode(ctx, driver, dbname)
	if err != nil {
		return nil, err
	}
	separatedio, err := SeparateIOType(ionodes)
	if err != nil {
		return nil, err
	}
	lgnodes, err := GetAllLogicNodes(ctx, driver, dbname)
	if err != nil {
		return nil, err
	}

	wnodes, err := GetAllWireNodes(ctx, driver, dbname)
	if err != nil {
		return nil, err
	}

	return &AllNode{
		Ios: separatedio,
		Lgs: lgnodes,
		Ws:  wnodes,
	}, nil
}

func SeparateIOType(allio map[string]*GetNeo4JIONode) (*IONodes, error) {
	innodes := make(map[string]*GetNeo4JIONode)
	outnodes := make(map[string]*GetNeo4JIONode)
	for key, io := range allio {
		switch io.ION.Type {
		case "IN":
			innodes[key] = io
		case "OUT":
			outnodes[key] = io
		default:
			return nil, fmt.Errorf("SeparateIOType Error %v", io.ION.Type)
		}
	}
	return &IONodes{
		In:  innodes,
		Out: outnodes,
	}, nil
}

// 親ノードを取得
// A <- (B) ,Get "B" Node
func Predecessors() {}

// 子ノードを取得
// (C) <- A ,Get "C" Node
func Successors() {}

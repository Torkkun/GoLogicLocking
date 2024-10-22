package veryosys

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// グラフとして表現するための部分

// Cells
type Cell struct {
	Type       string // gate or ff type
	Attributes struct {
		Src string
	}
}

// Connections
type Connections map[string][]*Connection // AttrSrc is Key
type Connection struct {
	Type     string // input or output
	BitNum   int
	PortName string
}

type NetName struct {
	Bits       []int
	Netname    string
	Attributes struct {
		Src string
	}
}

// Input Output InOut
type Port struct {
	Direction string
	BitNum    int //ポートの接続Bit
	BitWidth  int //ビット幅（すべてのBitが同じ場所に繋がっているわけではない）
	Name      string
}

// 全部のノードを取得
type AllNode struct {
	Ios  *IONodes
	Nets map[string]*NetName
	Cell map[string]*Cell
	// 増える可能性がある
}

// key is elementid
type IONodes struct {
	In  map[string]*Port
	Out map[string]*Port
	// InOut
}

func GetAllNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (*AllNode, error) {
	// In node Get
	// Out node Get
	// Cell node Get (premitive node)
	// Wire node Get
	return &AllNode{}, nil
}

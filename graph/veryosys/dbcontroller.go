package veryosys

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

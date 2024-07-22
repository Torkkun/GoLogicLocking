package veryosys

// Cells (Connection)
type Cell struct {
	Type        string           // gate or ff type
	Connections map[int]struct { // key is BitNum
		Type string // input or output
	}
}

// Wire
type NetName struct {
	BitNum    int
	Name      string
	WireGroup string
}

// Input Output InOut
type Port struct {
	Direction string
	BitNum    []int
	BitWidth  int
	Name      string
}

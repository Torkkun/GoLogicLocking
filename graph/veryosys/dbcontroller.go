package veryosys

// Cells
type Cell struct {
	Type        string
	Connections map[int]struct {
		Type string
	}
}

// return map[string]NetName
//
//	Key: ElementID, Value: NetName
//
// Wire
type NetName struct {
	BitNum    int
	Name      string
	WireGroup string
}

// Input Output InOut
type Port struct {
	Direction string
	BitNum    int
	BitWidth  int
	Name      string
}

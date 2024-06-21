package parser

import "goll/utils"

type prevType int

const (
	_ prevType = iota
	Lvalue
	Rvalue
	Gate
)

type Node struct {
	In1  string
	In2  string
	Gate utils.GateType
	Out  string
}

type LogicGate struct {
	GateType utils.GateType
}

type Decl struct {
	//name is key
	IOPorts map[string]IOPort
	Wires   map[string]Wire
}

type ParseResult struct {
	// relation
	Nodes        map[int]Node
	Declarations *Decl
	LogicGates   map[int]LogicGate
}

type IOPortType string

const (
	IN    IOPortType = "IN"
	OUT   IOPortType = "OUT"
	INOUT IOPortType = "INOUT"
)

type IOPort struct {
	Type IOPortType
	Name string
}

type Wire struct {
	Name string
	//todo
	//Width
}

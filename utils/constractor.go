package utils

type GateType string

const (
	Xor  GateType = "XOR"
	Xnor GateType = "XNOR"
	And  GateType = "AND"
	Nand GateType = "NAND"
	Or   GateType = "OR"
	Nor  GateType = "NOR"
	Not  GateType = "NOT"
)

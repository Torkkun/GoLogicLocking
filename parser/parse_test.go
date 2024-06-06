package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	parsed := NewParse("../testtxt/test.txt")
	fmt.Println(parsed.Nodes)
	fmt.Println(parsed.Declarations.IOPorts)
	fmt.Println(parsed.Declarations.Wires)
	fmt.Println(parsed.LogicGates)
}

package parser

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	parsed := NewParse()
	fmt.Println(parsed.Nodes)
	fmt.Println(parsed.Declaration.IOPorts)
	fmt.Println(parsed.Declaration.Wires)
}

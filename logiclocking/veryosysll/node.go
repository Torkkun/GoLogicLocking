package veryosysll

import (
	"context"
	"goll/graph/veryosys"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// MUX周りで何かしら変更する可能性があるかもしれないので一応
func AddLUTNode(tx neo4j.ExplicitTransaction, ctx context.Context, name string, muxwidth int) (string, error) {
	newlut := new(veryosys.LockLutNode)
	newlut.GateType = "LUT"
	newlut.LockType = "randomlut"
	newlut.Name = name

	elementId, err := veryosys.CreateRandomLutGateNodeTx(tx, ctx, newlut)
	if err != nil {
		return "", err
	}
	return elementId, nil
}

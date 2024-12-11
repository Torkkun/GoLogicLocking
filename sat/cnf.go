package sat

import (
	"context"
	"fmt"
	"goll/graph/veryosys"
	"slices"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// 回路の情報をCNFに変換する
// cnf()
func ConvertToCircuit(ctx context.Context, driver neo4j.DriverWithContext, dbname string, nodes []dbtype.Node) (*Formula, *IDPool, error) {
	//
	variable := NewIDPool(0, [][]int{})
	formula := NewCNF()
	for _, n := range nodes {
		variable.Id(n.GetElementId())
		if len(n.Labels) != 1 {
			return nil, nil, fmt.Errorf("too much labels: %d", len(n.Labels))
		}
		if n.Labels[0] != "Cell" {
			return nil, nil, fmt.Errorf("different labels type: %s", n.Labels[0])
		}
		gateType, ok := n.Props["type"].(string)
		if !ok {
			return nil, nil, fmt.Errorf("gateType is not a string")
		}
		// notやbufといった1ビットの場合はまだ決めてない

		// wireが間にあるので子がWireだった場合Cellにたどり着くまで子のノードを探索する必要がある
		sucs, err := veryosys.GetSuccessorNodes(ctx, driver, dbname, n.GetElementId())
		if err != nil {
			return nil, nil, err
		}

		switch gateType {
		case "$_NOT_":
			if sucs.Nodes != nil {
				// リストの一番後ろだけ削除
				f := sucs.Nodes[len(sucs.Nodes)]
				formula.Append([]int{variable.Id(n.GetElementId()), variable.Id(f.GetElementId())})
				formula.Append([]int{-variable.Id(n.GetElementId()), -variable.Id(f.GetElementId())})
			}
		case "$_BUF_":
			if sucs.Nodes != nil {
				// リストの一番後ろだけ削除
				f := sucs.Nodes[len(sucs.Nodes)]
				formula.Append([]int{variable.Id(n.GetElementId()), -variable.Id(f.GetElementId())})
				formula.Append([]int{-variable.Id(n.GetElementId()), variable.Id(f.GetElementId())})
			}
		case "$_AND_":
			clause := []int{variable.Id(n.GetElementId())}
			for _, f := range sucs.Nodes {
				formula.Append([]int{-variable.Id(n.GetElementId()), variable.Id(f.GetElementId())})
				clause = append(clause, -variable.Id(f.GetElementId()))
			}
			formula.Append(clause)
		case "$_NAND_":
			clause := []int{-variable.Id(n.GetElementId())}
			for _, f := range sucs.Nodes {
				formula.Append([]int{variable.Id(n.GetElementId()), variable.Id(f.GetElementId())})
				clause = append(clause, -variable.Id(f.GetElementId()))
			}
			formula.Append(clause)
		case "$_OR_":
			clause := []int{-variable.Id(n.GetElementId())}
			for _, f := range sucs.Nodes {
				formula.Append([]int{variable.Id(n.GetElementId()), -variable.Id(f.GetElementId())})
				clause = append(clause, variable.Id(f.GetElementId()))
			}
			formula.Append(clause)
		case "$_NOR_":
			clause := []int{variable.Id(n.GetElementId())}
			for _, f := range sucs.Nodes {
				formula.Append([]int{-variable.Id(n.GetElementId()), -variable.Id(f.GetElementId())})
				clause = append(clause, variable.Id(f.GetElementId()))
			}
			formula.Append(clause)
		case "$_XOR_", "$_XNOR_":
			var nets []string
			for _, node := range sucs.Nodes {
				nets = append(nets, node.GetElementId())
			}

			// xor gen
			xorClauses := func(a interface{}, b interface{}, c interface{}) {
				formula.Append([]int{-variable.Id(c), -variable.Id(b), -variable.Id(a)})
				formula.Append([]int{-variable.Id(c), variable.Id(b), variable.Id(a)})
				formula.Append([]int{variable.Id(c), -variable.Id(b), variable.Id(a)})
				formula.Append([]int{variable.Id(c), variable.Id(b), -variable.Id(a)})
			}

			for len(nets) > 2 {
				// create new net
				newnet := fmt.Sprintf("xor_%s_%s", nets[len(nets)-2], nets[len(nets)-1])
				variable.Id(newnet)

				// add sub xors
				xorClauses(nets[len(nets)-2], nets[len(nets)-1], newnet)

				// remove last 2 nets
				nets = nets[:len(nets)-2]

				// insert
				nets = slices.Insert(nets, 0, newnet)
			}
			//add final xor
			if gateType == "$_XOR_" {
				xorClauses(nets[len(nets)-2], nets[len(nets)-1], n.GetElementId())
			} else {
				tmpname := fmt.Sprintf("xor_inv_{%s}", n.GetElementId())
				variable.Id(tmpname)
				xorClauses(nets[len(nets)-2], nets[len(nets)-1], tmpname)
				formula.Append([]int{variable.Id(n.GetElementId()), variable.Id(tmpname)})
				formula.Append([]int{-variable.Id(n.GetElementId()), -variable.Id(tmpname)})
			}
		case "0":
			formula.Append([]int{-variable.Id(n.GetElementId())})
		case "1":
			formula.Append([]int{variable.Id(n.GetElementId())})
		default:
			return nil, nil, fmt.Errorf("unknown gate type %v", gateType)
		}
	}
	return formula, variable, nil
}

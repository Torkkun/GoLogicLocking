package sat

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/dbtype"
)

// 一時的なもの仮置き
type Solver struct {
}

func Solve(ctx context.Context, driver neo4j.DriverWithContext, dbname string, nodes []dbtype.Node, assumptions map[string]bool) (map[string]bool, error) {
	// ライブラリを使う場合は、ソルバインスタンスを生成今はとりあえずCMD実行予定
	//solver, variables := construct_solver(ctx, driver, dbname, nodes, assumptions)

	// CNF

	ConvertToCircuit(ctx, driver, dbname, nodes)

	return nil, nil
}

// コマンドプロンプトから実行
func ExecCmd() {

}

// ライブラリからはするかもしれない

func construct_solver(ctx context.Context, driver neo4j.DriverWithContext, dbname string, nodes []dbtype.Node, assumptions map[string]bool) (*Solver, *IDPool, error) {
	//一時的にformulaをアウト
	_, variables, err := ConvertToCircuit(ctx, driver, dbname, nodes)
	if err != nil {
		return nil, nil, err
	}
	return &Solver{}, variables, nil
}

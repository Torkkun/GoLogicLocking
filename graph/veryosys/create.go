package veryosys

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateWireNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, netname *NetName) (string, error) {

	return "", nil
}

func CreateCellNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, cell *Cell) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (cell:Cell {type: $type})`,
		map[string]any{
			"type": cell.Type,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return "", err
	}
	if len(result.Records) > 1 {
		return "", fmt.Errorf("to much record error")

	}
	record := result.Records[0]
	newcell, ok := record.Get("cell")
	if !ok {
		err = fmt.Errorf("GetCell Error")
		return "", err
	}
	cellnode := newcell.(neo4j.Node)

	return cellnode.GetElementId(), nil
}

package veryosys

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Node struct {
	Type      string
	ElementId string
}

type ConnectionPair struct {
	Predecessor Node
	Successor   Node
}

func CellConnection(ctx context.Context, driver neo4j.DriverWithContext, dbname string, conn *ConnectionPair) error {
	query := fmt.Sprintf(`
		MATCH (pre:%s), (suc:%s)
		WHERE elementId(pre)=$pre_element_id AND elementId(suc)=$suc_element_id
		MERGE (pre)-[:%sto%s]->(suc)`,
		conn.Predecessor.Type,
		conn.Successor.Type,
		conn.Predecessor.Type,
		conn.Successor.Type,
	)

	_, err := neo4j.ExecuteQuery(ctx, driver,
		query,
		map[string]any{
			"pre_element_id": conn.Predecessor.ElementId,
			"suc_element_id": conn.Successor.ElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE Cell Connection Error:%v", err)
		return err
	}
	return nil
}

func CellConnectionTx(tx neo4j.ExplicitTransaction, ctx context.Context, conn *ConnectionPair) error {
	query := fmt.Sprintf(`
		MATCH (pre:%s), (suc:%s)
		WHERE elementId(pre)=$pre_element_id AND elementId(suc)=$suc_element_id
		MERGE (pre)-[:%sto%s]->(suc)`,
		conn.Predecessor.Type,
		conn.Successor.Type,
		conn.Predecessor.Type,
		conn.Successor.Type,
	)

	_, err := tx.Run(ctx,
		query,
		map[string]any{
			"pre_element_id": conn.Predecessor.ElementId,
			"suc_element_id": conn.Successor.ElementId,
		})
	if err != nil {
		err = fmt.Errorf("MERGE Cell Connection Error:%v", err)
		return err
	}
	return nil
}

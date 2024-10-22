package veryosys

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func DeleteConnection(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH()-[r]->()
		WHERE elementId(r)=$element_id
		DELETE r`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DELETE Relationship Error:%v", err)
		return err
	}
	return nil
}

func DeleteConnectionTx(tx neo4j.ExplicitTransaction, ctx context.Context, elementId string) error {
	_, err := tx.Run(ctx,
		`MATCH()-[r]->()
		WHERE elementId(r)=$element_id
		DELETE r`,
		map[string]any{
			"element_id": elementId,
		})
	if err != nil {
		err = fmt.Errorf("DELETE Relationship Error:%v", err)
		return err
	}
	return nil
}

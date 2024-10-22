package veryosys

import (
	"context"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Transaction(ctx context.Context, driver neo4j.DriverWithContext, dbname string, txFunc func(neo4j.ExplicitTransaction) (interface{}, error)) (interface{}, error) {
	session := driver.NewSession(ctx, neo4j.SessionConfig{DatabaseName: dbname})
	defer session.Close(ctx)
	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(ctx); err != nil {
				panic(p)
			}
		} else if err != nil {
			if err := tx.Rollback(ctx); err != nil {
				log.Println(err)
			}
			return
		}
		err = tx.Commit(ctx)
		if err != nil {
			log.Println(err)
		}
	}()
	return txFunc(tx)
}

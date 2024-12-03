package sat

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Solver(ctx context.Context, driver neo4j.DriverWithContext, dbname string, assumptions map[string]bool) (map[string]bool, error) {

}

package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Delete Relationship
func (gi *GateIO) DeleteRelationGatetoIO(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (io:IO {type: $io_type, name: $io_name})<-[r:LGtoIO]-(g:Gate {type: $g_type, at: $g_at})
		DELETE r`,
		map[string]any{
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
			"g_type":  gi.Gate.GateType,
			"g_at":    gi.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DELETE GatetoIO Relation Error:%v", err)
		return err
	}
	return nil
}

func (gi *GateIO) DeleteRelationIOtoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at})<-[r:IOtoLG]-(io:IO {type: $io_type, name: $io_name})
		DELETE r`,
		map[string]any{
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
			"g_type":  gi.Gate.GateType,
			"g_at":    gi.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DELETE IOtoGate Relation Error:%v", err)
		return err
	}
	return nil
}

func DeleteRelationIOtoGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (:Gate)<-[r:IOtoLG]-(:IO)
		WHERE elementId(r)=$element_id
		DELETE r`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DELETE IOtoGate Relation Error:%v", err)
		return err
	}
	return nil
}

func (gw *GateWire) DeleteRelationWiretoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at})<-[r:WiretoLG]-(w:Wire {name: $w_name})
		DELETE r`,
		map[string]any{
			"g_type": gw.Gate.GateType,
			"g_at":   gw.Gate.At,
			"w_name": gw.Wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("DELETE WiretoGate Relation Error:%v", err)
		return err
	}
	return nil
}

func DeleteRelationWiretoGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname, elementId string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (:Gate)<-[r:WiretoLG]-(:Wire)
		WHERE elementId(r)=$element_id
		DELETE r`,
		map[string]any{
			"element_id": elementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("DELETE WiretoGate Relation Error:%v", err)
		return err
	}
	return nil
}

func (gw *GateWire) DeleteRelationGatetoWire(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (w:Wire {name: $w_name})<-[r:LGtoWire]-(g:Gate {type: $g_type, at: $g_at})
		DELETE r `,
		map[string]any{
			"w_name": gw.Wire.Name,
			"g_type": gw.Gate.GateType,
			"g_at":   gw.Gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("DELETE GatetoWire Relation Error:%v", err)
		return err
	}
	return nil
}

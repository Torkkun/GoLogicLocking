package verpyverilog

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type LockGateIO struct {
	Gate          *LockGateNode
	GateElementId string
	Io            *IONode
	IoElementId   string
}

/* func (lgi *LockGateIO) LLGatetoIO(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (io:IO {type: $io_type, name: $io_name}), (g:LLGate {type: $g_type, locktype: $ll_type, name: $name})
		MERGE (io)<-[:LLGtoIO]-(g)
		`,
		map[string]any{
			"io_type": lgi.Io.Type,
			"io_name": lgi.Io.Name,
			"g_type":  lgi.Gate.GateType,
			"ll_type": lgi.Gate.LockType,
			"name":    lgi.Gate.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
} */

func (lgi *LockGateIO) LLGatetoIOByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (io:IO), (llg:LLGate)
		WHERE elementId(io)=$io_element_id AND elementId(llg)=$llg_element_id
		MERGE (io)<-[:LLGtoIO]-(llg)
		`,
		map[string]any{
			"io_element_id":  lgi.IoElementId,
			"llg_element_id": lgi.GateElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
}

/* type LLGateGate struct {
	LLGN          *LockGateNode
	LLGNElementId string
	LGN           *LogicGateNode
	LGNElementId  string
} */

// Gate <- LLGate
// 間にワイヤーを作成しなければならない
/* func (llgg *LLGateGate) LLGatetoGateElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string, num int) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate), (lg:LLGate)
		WHERE elementId(lg)=$lg_element_id AND elementId(g)=$g_element_id
		MERGE (g)<-[:WiretoLG]-(:Wire {name: $w_name})<-[:LLGtoWire]-(lg)`,
		//MERGE (g)<-[:LLGtoG]-(lg)`,
		map[string]any{
			"lg_element_id": llgg.LLGNElementId,
			"g_element_id":  llgg.LGNElementId,
			"w_name":        fmt.Sprintf("llw_%d", num),
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoGate Error:%v", err)
		return err
	}
	return nil
} */

type LockGateWire struct {
	Gate          *LockGateNode
	GateElementId string
	Wire          *WireNode
	WireElementId string
}

// LLGate <- Wire
func (lgw *LockGateWire) WiretoLLGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (lg:LLGate), (w:Wire)
		WHERE elementId(lg)=$lg_element_id AND elementId(w)=$w_element_id
		MERGE (lg)<-[:WiretoLLG]-(w)
		`,
		map[string]any{
			"lg_element_id": lgw.GateElementId,
			"w_element_id":  lgw.WireElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
}

// Wire <- LLGate
func (lgw *LockGateWire) LLGatetoWireByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (w:Wire), (lg:LLGate)
		WHERE  elementId(w)=$w_element_id AND elementId(lg)=$lg_element_id
		MERGE (w)<-[:LLGtoWire]-(lg)
		`,
		map[string]any{
			"w_element_id":  lgw.WireElementId,
			"lg_element_id": lgw.GateElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE LLGatetoIO Error:%v", err)
		return err
	}
	return nil
}

// Create Relationship
// 後で良い方法を考える OGM（object graph mapping）
// IO(IN) <- Gate

type GateIO struct {
	Gate          *LogicGateNode
	GateElementId string
	Io            *IONode
	IoElementId   string
}

/* func (gi *GateIO) GatetoIO(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver, `
		MATCH (io:IO {type: $io_type, name: $io_name}), (g:Gate {type: $g_type, at: $g_at})
		MERGE (io)<-[:LGtoIO]-(g)
		`,
		map[string]any{
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
			"g_type":  gi.Gate.GateType,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE GatetoIO Error:%v", err)
		return err
	}
	return nil
} */

func (gi *GateIO) GatetoIOByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (io:IO), (g:Gate)
	WHERE elementId(io)=$io_element_id AND elementId(g)=$g_element_id
	MERGE (io)<-[:LGtoIO]-(g)`,
		map[string]any{
			"io_element_id": gi.IoElementId,
			"g_element_id":  gi.GateElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE GatetoIO Error:%v", err)
		return err
	}
	return nil
}

// Gate <- IO(OUT)
/* func (gi *GateIO) IOtoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type}), (io:IO {type: $io_type, name: $io_name})
		MERGE (g)<-[:IOtoLG]-(io)`,
		map[string]any{
			"g_type":  gi.Gate.GateType,
			"io_type": gi.Io.Type,
			"io_name": gi.Io.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE IOtoGate Error:%v", err)
		return err
	}
	return nil
} */

func (gi *GateIO) IOtoGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate), (io:IO)
		WHERE elementId(g)=$g_element_id AND elementId(io)=$io_element_id
		MERGE (g)<-[:IOtoLG]-(io)`,
		map[string]any{
			"g_element_id":  gi.GateElementId,
			"io_element_id": gi.IoElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE IOtoGate Error:%v", err)
		return err
	}
	return nil
}

// LockingGate <- IO(OUT)
func (lgi *LockGateIO) IOtoLLGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (lg:LLGate), (io:IO)
		WHERE elementId(lg)=$lg_element_id AND elementId(io)=$io_element_id
		MERGE (lg)<-[:IOtoLLG]-(io)`,
		map[string]any{
			"lg_element_id": lgi.GateElementId,
			"io_element_id": lgi.IoElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE IOtoLLGate Error:%v", err)
		return err
	}
	return nil
}

type GateWire struct {
	Gate          *LogicGateNode
	GateElementId string
	Wire          *WireNode
	WireElementId string
}

// Gate <- Wire
/* func (gw *GateWire) WiretoGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate {type: $g_type, at: $g_at}), (w:Wire {name: $w_name})
		MERGE (g)<-[:WiretoLG]-(w)`,
		map[string]any{
			"g_type": gw.Gate.GateType,
			"w_name": gw.Wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE WiretoGate Error:%v", err)
		return err
	}
	return nil
} */

func (gw *GateWire) WiretoGateByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (g:Gate), (w:Wire)
		WHERE elementId(g)=$g_element_id AND elementId(w)=$w_element_id
		MERGE (g)<-[:WiretoLG]-(w)`,
		map[string]any{
			"g_element_id": gw.GateElementId,
			"w_element_id": gw.WireElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE WiretoGate Error:%v", err)
		return err
	}
	return nil
}

// Wire <- Gate
/* func (gw *GateWire) GatetoWire(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (w:Wire {name: $w_name}), (g:Gate {type: $g_type, at: $g_at})
		MERGE (w)<-[:LGtoWire]-(g)`,
		map[string]any{
			"w_name": gw.Wire.Name,
			"g_type": gw.Gate.GateType,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("MERGE GatetoWire Error:%v", err)
		return err
	}
	return nil
} */

func (gw *GateWire) GatetoWireByElementId(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (w:Wire), (g:Gate)
		WHERE elementId(w)=$w_element_id AND elementId(g)=$g_element_id
		MERGE (w)<-[:LGtoWire]-(g)`,
		map[string]any{
			"w_element_id": gw.WireElementId,
			"g_element_id": gw.GateElementId,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("MERGE GatetoWire Error:%v", err)
		return err
	}
	return nil
}

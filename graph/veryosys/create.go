package veryosys

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func CreateInOutNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, port *Port) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (io:IO {direction: $direction, name:$name, bitnum: $bitnum, bitwidth: $bitwidth})
		RETURN io`,
		map[string]any{
			"direction": port.Direction,
			"name":      port.Name,
			"bitnum":    port.BitNum,
			"bitwidth":  port.BitWidth,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return "", err
	}
	record := result.Records[0]
	resio, ok := record.Get("io")
	if !ok {
		err = fmt.Errorf("NotfoundCreated IO Node")
		return "", err
	}
	n := resio.(neo4j.Node)
	return n.GetElementId(), nil
}

func CreateInOutNodeTx(tx neo4j.ExplicitTransaction, ctx context.Context, port *Port) (string, error) {
	result, err := tx.Run(ctx,
		`CREATE (io:IO {direction: $direction, name:$name, bitnum: $bitnum, bitwidth: $bitwidth})
		RETURN io`,
		map[string]any{
			"direction": port.Direction,
			"name":      port.Name,
			"bitnum":    port.BitNum,
			"bitwidth":  port.BitWidth,
		})
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return "", err
	}
	record := result.Record()
	resio, ok := record.Get("io")
	if !ok {
		err = fmt.Errorf("NotfoundCreated IO Node")
		return "", err
	}
	n := resio.(neo4j.Node)
	return n.GetElementId(), nil
}

type DBNetName struct {
	BitNum     int
	Netname    string
	Attributes struct {
		Src string
	}
}

func CreateWireNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, netname *DBNetName) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (w:Wire {bitnum: $bitnum, netname: $netname, attrsrc: $attrsrc})
			RETURN w`,
		map[string]any{
			"bitnum":  netname.BitNum,
			"netname": netname.Netname,
			"attrsrc": netname.Attributes.Src,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("CreateNetWireNode Error:%v", err)
		return "", err
	}
	record := result.Records[0]
	resw, ok := record.Get("w")
	if !ok {
		err = fmt.Errorf("NotfoundCreated Wire Node")
		return "", err
	}
	w := resw.(neo4j.Node)

	return w.GetElementId(), nil
}

// GateLevel Node
func CreateCellNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, cell *Cell) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (cell:Cell {type: $type, attrsrc: $attrsrc})
		RETURN cell`,
		map[string]any{
			"type":    cell.Type,
			"attrsrc": cell.Attributes.Src,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateCellNode Error:%v", err)
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

type LockGateNode struct {
	GateType string
	LockType string
}

func CreateRandomLLGateNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string, llnode *LockGateNode) (string, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (llcell:LLCell {gatetype: $gate_type, locktype: $lock_type})
		RETURN llcell`,
		map[string]any{
			"gate_type": llnode.GateType,
			"lock_type": llnode.LockType,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateLLCellNode Error:%v", err)
		return "", err
	}
	record := result.Records[0]
	llcell, ok := record.Get("llcell")
	if !ok {
		err = fmt.Errorf("GetCell Error")
		return "", err
	}
	llcellnode := llcell.(neo4j.Node)

	return llcellnode.GetElementId(), nil
}

func CreateRandomLLGateNodeTx(tx neo4j.ExplicitTransaction, ctx context.Context, llnode *LockGateNode) (string, error) {
	result, err := tx.Run(ctx,
		`CREATE (llcell:LLCell {gatetype: $gate_type, locktype: $lock_type})
		RETURN llcell`,
		map[string]any{
			"gate_type": llnode.GateType,
			"lock_type": llnode.LockType,
		},
	)
	if err != nil {
		err = fmt.Errorf("CreateLLCellNode Error:%v", err)
		return "", err
	}
	record := result.Record()
	llcell, ok := record.Get("llcell")
	if !ok {
		err = fmt.Errorf("GetCell Error")
		return "", err
	}
	llcellnode := llcell.(neo4j.Node)

	return llcellnode.GetElementId(), nil
}

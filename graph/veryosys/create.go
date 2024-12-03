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
	record, err := result.Single(ctx)
	if err != nil {
		return "", err
	}
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

func CreateWireNodeTx(tx neo4j.ExplicitTransaction, ctx context.Context, netname *DBNetName) (string, error) {
	result, err := tx.Run(ctx,
		`CREATE (w:Wire {bitnum: $bitnum, netname: $netname, attrsrc: $attrsrc})
			RETURN w`,
		map[string]any{
			"bitnum":  netname.BitNum,
			"netname": netname.Netname,
			"attrsrc": netname.Attributes.Src,
		})
	if err != nil {
		err = fmt.Errorf("CreateNetWireNode Error:%v", err)
		return "", err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return "", err
	}
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

func CreateCellNodeTx(tx neo4j.ExplicitTransaction, ctx context.Context, cell *Cell) (string, error) {
	result, err := tx.Run(ctx,
		`CREATE (cell:Cell {type: $type, attrsrc: $attrsrc})
		RETURN cell`,
		map[string]any{
			"type":    cell.Type,
			"attrsrc": cell.Attributes.Src,
		}) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateCellNode Error:%v", err)
		return "", err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return "", err
	}
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
	record, err := result.Single(ctx)
	if err != nil {
		return "", err
	}
	llcell, ok := record.Get("llcell")
	if !ok {
		err = fmt.Errorf("GetCell Error")
		return "", err
	}
	llcellnode := llcell.(neo4j.Node)

	return llcellnode.GetElementId(), nil
}

type LockLutNode struct {
	GateType string
	LockType string
	Name     string
}

func CreateRandomLutGateNodeTx(tx neo4j.ExplicitTransaction, ctx context.Context, llnode *LockLutNode) (string, error) {
	result, err := tx.Run(ctx,
		`CREATE (llLut:LLLut {gatetype: $gate_type, locktype: $lock_type})
		RETURN llLut`,
		map[string]any{
			"gate_type": llnode.GateType,
			"lock_type": llnode.LockType,
			"name":      llnode.Name,
		},
	)
	if err != nil {
		err = fmt.Errorf("CreateLLLutNode Error:%v", err)
		return "", err
	}
	record, err := result.Single(ctx)
	if err != nil {
		return "", err
	}
	llcell, ok := record.Get("llLut")
	if !ok {
		err = fmt.Errorf("GetLut Error")
		return "", err
	}
	lllutnode := llcell.(neo4j.Node)

	return lllutnode.GetElementId(), nil
}

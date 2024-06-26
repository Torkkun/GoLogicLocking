package verpyverilog

type Relation struct {
	Identity           int64
	Type               string
	Properties         interface{}
	ElementId          string
	StartNodeElementId string
	EndNodeElementId   string
}

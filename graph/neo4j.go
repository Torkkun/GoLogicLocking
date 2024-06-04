package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type graphDB struct {
	Driver neo4j.DriverWithContext
	DBname string
}

// test
func NewDriver() *graphDB {
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	return &graphDB{Driver: driver}
}

// setting Jsonから取得するようにしたい
func SelectDriver(dbselect string) *graphDB {
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "password"
	var dbport string
	switch dbselect {
	case "origin":
		dbport = "7688"
	case "ll":
		dbport = "7689"
	default:
		dbport = ""
	}
	if dbport == "" {
		panic(fmt.Sprintf("ERROR: dbname is null = %s", dbselect))
	}
	driver, err := neo4j.NewDriverWithContext(
		fmt.Sprintf(dbUri+":%s", dbport),
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	return &graphDB{Driver: driver, DBname: "neo4j"}
}

type IONode struct {
	Type string
	Name string
}

func (io *IONode) CreateInOutNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:IO {type: $type, name:$name})`,
		map[string]any{
			"type": io.Type,
			"name": io.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択 後で
	if err != nil {
		err = fmt.Errorf("CreateInOutNode Error:%v", err)
		return err
	}
	return nil
}

type GetNeo4JIONode struct {
	ION       *IONode
	Id        int64
	ElementId string
}

func GetAllIONode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) ([]*GetNeo4JIONode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(io:IO) RETURN io`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllIONodes Error:%v", err)
		return nil, err
	}
	var ionodes []*GetNeo4JIONode
	for _, record := range result.Records {
		io, ok := record.Get("io")
		if !ok {
			err = fmt.Errorf("GetAllIONode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		ionodes = append(ionodes,
			&GetNeo4JIONode{
				ION: &IONode{
					Type: tmp.Props["type"].(string),
					Name: tmp.Props["name"].(string),
				},
				ElementId: tmp.GetElementId(),
				Id:        tmp.GetId(),
			},
		)
	}
	return ionodes, nil
}

type LogicGateNode struct {
	GateType string
	At       int
}

func (gate *LogicGateNode) CreateLogicGateNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Gate {type: $type, at: $at})`,
		map[string]any{
			"type": gate.GateType,
			"at":   gate.At,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateLogicGateNode Error:%v", err)
		return err
	}
	return nil
}

type GetNeo4JLogicNode struct {
	LGN       *LogicGateNode
	Id        int64
	ElementId string
}

func GetAllLogicNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) ([]*GetNeo4JLogicNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(g:Gate) RETURN g`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllLogicNodes Error:%v", err)
		return nil, err
	}
	var lgnodes []*GetNeo4JLogicNode
	for _, record := range result.Records {
		io, ok := record.Get("g")
		if !ok {
			err = fmt.Errorf("GetAllLogicNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		lgnodes = append(lgnodes,
			&GetNeo4JLogicNode{
				LGN: &LogicGateNode{
					GateType: tmp.Props["type"].(string),
					At:       int(tmp.Props["at"].(int64)),
				},
				ElementId: tmp.GetElementId(),
				Id:        tmp.GetId(),
			})
	}
	return lgnodes, nil
}

type WireNode struct {
	Name string
}

func (wire *WireNode) CreateWireNode(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`CREATE (:Wire {name: $name})`,
		map[string]any{
			"name": wire.Name,
		},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname)) //DBの選択
	if err != nil {
		err = fmt.Errorf("CreateWireNode Error:%v", err)
		return err
	}
	return nil
}

type GetNeo4JWireNode struct {
	WN        *WireNode
	Id        int64
	ElementId string
}

func GetAllWireNodes(ctx context.Context, driver neo4j.DriverWithContext, dbname string) ([]*GetNeo4JWireNode, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH(w:Wire) RETURN w`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("GetAllNodes Error:%v", err)
		return nil, err
	}
	var wnodes []*GetNeo4JWireNode
	for _, record := range result.Records {
		io, ok := record.Get("w")
		if !ok {
			err = fmt.Errorf("GetAllWireNode Error")
			return nil, err
		}
		tmp := io.(neo4j.Node)
		wnodes = append(wnodes, &GetNeo4JWireNode{
			WN: &WireNode{
				Name: tmp.Props["name"].(string),
			},
			ElementId: tmp.GetElementId(),
			Id:        tmp.GetId(),
		})
	}
	return wnodes, nil
}

// 複数のユーザーデータベースは作れないので全部消すしかないため用意
func DBtableAllClear(ctx context.Context, driver neo4j.DriverWithContext, dbname string) error {
	_, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (n) DETACH DELETE n`,
		nil,
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("DBtableAllClear Error")
		return err
	}
	return nil
}

func CountGate(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (int, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (n:Gate) RETURN count(n) as count`,
		map[string]any{},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("CountGateNode Error:%v", err)
		return 0, err
	}
	var tmp int64
	if len(result.Records) != 1 {
		err = fmt.Errorf("CountGateNode query is not invalid")
		return 0, err
	} else {
		record := result.Records[0]
		rawValue, ok := record.Get("count")
		if !ok {
			err = fmt.Errorf("CountGateNode Error")
			return 0, err
		}
		tmp = rawValue.(int64)
	}
	intnum := int(tmp)
	return intnum, nil
}

func CountOUT(ctx context.Context, driver neo4j.DriverWithContext, dbname string) (int, error) {
	result, err := neo4j.ExecuteQuery(ctx, driver,
		`MATCH (n:IO {type:"OUT"}) RETURN count(n) as count`,
		map[string]any{},
		neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(dbname))
	if err != nil {
		err = fmt.Errorf("CountIONode Error:%v", err)
		return 0, err
	}
	var tmp int64
	if len(result.Records) != 1 {
		err = fmt.Errorf("CountIONode query is not invalid")
		return 0, err
	} else {
		record := result.Records[0]
		rawValue, ok := record.Get("count")
		if !ok {
			err = fmt.Errorf("CountIONode Error")
			return 0, err
		}
		tmp = rawValue.(int64)
	}
	intnum := int(tmp)
	return intnum, nil
}

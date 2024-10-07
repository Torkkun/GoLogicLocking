package graph

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphDB struct {
	Driver neo4j.DriverWithContext
	DBname string
}

// test
func NewDriver() *GraphDB {
	dbUri := "neo4j://localhost"
	dbUser := "neo4j"
	dbPassword := "secretgraph"
	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(dbUser, dbPassword, ""))
	if err != nil {
		panic(err)
	}
	return &GraphDB{Driver: driver, DBname: "neo4j"}
}

// setting Jsonから取得するようにしたい
func SelectDriver(dbselect string) *GraphDB {
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
	return &GraphDB{Driver: driver, DBname: "neo4j"}
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

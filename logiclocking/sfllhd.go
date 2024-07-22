package logiclocking

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func SfllHd(ctx context.Context, driver neo4j.DriverWithContext, dbname string, keylen int) (map[string]bool, error) {
	// count circuitを用意
	// 元の回路の入力数がキー幅を下回るとエラー
	// circuit input < key width
	// ターゲットとするアウトプットの選択がなければランダムに選択する
	//
}

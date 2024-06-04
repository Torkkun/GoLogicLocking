package logiclocking

// TODO
func XorLock() {
	// db側でコピー（Cyperクエリのバックアップクエリ）

	// ノード数 - アウトプットの数でロック用のゲートを作成
	//  - （ランダムにXORとXNORのゲート情報を作成していると思われる）

	// 上のゲートをひとつづつ取り出し、キーリスト（map）にてTrueまたはFalseを選択
	//  - circuitgraphのchoicesとsample関数をチラ見する

	// TRUEまたはFALSEがゲートに設定されるが、
	//  - TRUEならXNOR、FALSEならXORとゲートタイプが決定される

	// ゲートの追加とインプットの追加を行う
	//　もともとのロジックゲートのOUTと
}

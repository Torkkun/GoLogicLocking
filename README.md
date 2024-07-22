# GoLogicLocking

[] metrics
[] sat solver & cnf convertor
[] multi bit対応
[] いくつかの難読化

## GrapGB（neo4j）
DBを二つ建てて一つをオリジナル、もう一つをLogicLocking用のDBにする
オリジナルにて最初にグラフをＹｏｓｙｓにて論理合成したあとのＪｓｏｎ形式のネットリストからグラフを構成
構成後、neo4jからcsvとしてexport

json, csvをnosqlDBの何かにバックアップとして保存、管理
ロジックロックの展開の際にcsvを読み込むロッキングした回路情報を保存したければ同様にバックアップcsvを保存

このように管理する方法にする？

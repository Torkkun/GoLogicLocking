# Json構文
https://yosyshq.readthedocs.io/projects/yosys/en/latest/cmd/write_json.html

# 内部セルライブラリ
https://blog.eowyn.net/yosys/CHAPTER_CellLib.html

組み合わせ回路に対しては、LogicLockingを行えるが、順序回路は無視したいためそれらを識別する

元のverilog 論理合成前
↓　yosysにて論理合成
ゲートレベル verilog

NetNames 使用されるワイヤーの情報
hide_nameは自動的に作成されるものが1、InOutなどverilogで記述したものは0に設定される。

# Cell部分のマニュアル
## Connected部分のbit Vectorは、整数の場合や文字列の場合がある
https://yosyshq.readthedocs.io/projects/yosys/en/latest/cmd/write_json.html
ジェネリクスを使用するのがよさそうだがどのような場合においてこのようになるパターンがあるのかわからない

わかりにくすぎる

一旦保留
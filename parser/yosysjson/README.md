# Json構文
https://yosyshq.readthedocs.io/projects/yosys/en/latest/cmd/write_json.html

# 内部セルライブラリ
https://blog.eowyn.net/yosys/CHAPTER_CellLib.html

組み合わせ回路に対しては、LogicLockingを行えるが、順序回路は無視したいためそれらを識別する

元のverilog 論理合成前
↓　yosysにて論理合成
ゲートレベル verilog

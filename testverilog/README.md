# dockerを使用する場合
```shell
sudo docker build -t testcg .
sudo docker run --mount type=bind,source=./workspace,target=/work/workspace --name testcg -it circuitgc-test /bin/bash
```
python3 main.pyを実行する

# yosys環境変数
```shell
source ~/oss-cad-suite/environment
```
# yosys
```shell
$ yosys
```

Verilog フロントエンドを使用して設計を読み取って詳しく説明します。
```shell
$ read -sv path/to/verilog/??.v
$ hierarchy -top ??
```
Yosys が内部的に使用する RTLIL 形式でデザインをコンソールに書き込みます。
```shell
$ write_rtlil
```
プロセス (alwaysブロック) をネットリスト要素に変換し、いくつかの簡単な最適化を実行します。
```shell
$ proc; opt
```
次を使用してデザイン ネットリストを表示しますxdot。
```shell
$ show
```
ポストスクリプトビューアと同じものを使用しますgv:
```shell
$ show -format ps -viewer gv
```
ネットリストをゲート ロジックに変換し、いくつかの簡単な最適化を実行します。
```shell
$ techmap; opt
```
デザイン ネットリストを新しい Verilog ファイルに書き込みます。
```shell
$ write_verilog synth.v
```
##

synthesisscript
```shell
$ yosys synth.ys
```
おそらくPythonのやつは配列とかは対応してなさそう（issueからマルチバイト対応を行うやつを試し中）

SATソルバとしてcadicalというものを使用
https://github.com/arminbiere/cadical

wslの~/cadicalにビルド済み ver2.1.0
~/cadical/build -t 120 cadicaltest/{適当な例}　にて実行できる
windowsからは
wsl.exe ~/cadical/build -t 120 cadicaltest/{適当な例}

ライブラリとして実行するには以下のようにただしWindowsでは無理かも
// SATSolverのライブラリ例としてPythonとRustがある
https://github.com/pysathq/pysat
https://github.com/chrjabs/rustsat


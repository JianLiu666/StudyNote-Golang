# Benchmark 壓力測試

執行所有的 Benchmark 測試並產生對應的 Profiling 結果檔案

 - -bench: 測試函式名稱 (預設=.表示全部都測試)
 - -benchtime: 測試次數/時間 (e.g. 1000x=1000次, 10s=10秒)
 - -cpu: 測試時使用的cpu數目
 - -cpuprofile: cpu profiling 數據要保存的位置
 - -memprofile: memory profiling 數據要保存的位置

```sh
go test -bench={funcName} -benchtime=10000000x -cpu=1 -cpuprofile prof.cpu -memprofile prof.mem
```

# pprof web UI

直接讀取 Profiling 檔案以網頁方式呈現

```sh
go tool pprof -http=":8081" pprof.test prof.cpu
go tool pprof -http=":8081" pprof.test prof.mem
```

# pprof commands

使用 pprof 工具開始進行效能分析

```sh
go tool pprof pprof.test prof.cpu
go tool pprof pprof.test prof.mem
```

### top N

列出前 N 名最吃效能的函式列表 (i.e. 最耗時、使用最多記憶體)

 - 前兩列表示函數佔用的效能與百分比
 - 第三列是當前所有函數累加所佔用的百分比
 - 第四、五兩列表示當前函數與所呼叫子函數的佔用效能與百分比

### web

以瀏覽器打開函數調用的效能分析圖

### list [benchName]

列出函數程式碼以及對應的效能分析結果

# Reference

1. 安裝 graphviz 使用效能分析圖與火焰圖

```sh
sudo apt-get install -y graphviz
```

2. [使用 pprof 和火焰圖調試 golang 應用](https://cizixs.com/2017/09/11/profiling-golang-program/)

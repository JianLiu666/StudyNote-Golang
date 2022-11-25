## Benchmark 壓力測試

參數筆記：

 - `-bench regexp` : 壓力測試函式名稱 (預設 `.` 表示全部都測試)
 - `-run regexp` : 單元測試函數名稱 (預設 `.` 表示全部都測試)
 - `-benchtime t` : 測試次數或時間 (e.g. 1000x=1000次, 10s=10秒)
 - `-cpu 1,2,4` : 測試時使用的cpu數目
 - `-benchmem` : 
 - `-cpuprofile cpu.out` : 輸出 cpu profile
 - `-memprofile mem.out` : 輸出 memory profile

更多參數說明可以參考：

```sh
go help testflag
```

範例：執行當前路徑底下所有的測試函式並產生對應的 cpu 與 memory profiles

```sh
go test -bench=. -run=none -benchtime=10000000x -cpu=1 -benchmem -cpuprofile cpu.profile -memprofile mem.profile
```

---

## pprof web UI

讀取指定的 profile，並以網頁呈現

```sh
go tool pprof -http=":8081" cpu.profile
go tool pprof -http=":8081" mem.profile
```

---

## pprof commands

使用 pprof CLI 進行效能分析

```sh
go tool pprof pprof.test cpu.profile
go tool pprof pprof.test mem.profile
```

### top [N]

列出前 N 名最吃效能的函式列表 (i.e. 最耗時、使用最多記憶體)

 - 前兩列表示函數佔用的效能與百分比
 - 第三列是當前所有函數累加所佔用的百分比
 - 第四、五兩列表示當前函數與所呼叫子函數的佔用效能與百分比

### web

以瀏覽器打開函數調用的效能分析圖

### list [benchName]

列出函數程式碼以及對應的效能分析結果

---

## References

1. 安裝 graphviz 使用效能分析圖與火焰圖

```sh
sudo apt-get install -y graphviz
```

2. [使用 pprof 和火焰圖調試 golang 應用](https://cizixs.com/2017/09/11/profiling-golang-program/)

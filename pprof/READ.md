# 產生測試檔案

執行 Benchmark 測試並產生對應的測試檔案
> go test . -bench . -cpuprofile prof.cpu

讀取測試檔案
> go tool pprof pprof.test prof.cpu

# pprof 介面

觀察執行最慢的前10名動作
> top10

以瀏覽器打開執行效能分析圖
> web

# Reference

1. 安裝 graphviz 使用效能分析圖與火焰圖
> sudo apt-get install -y graphviz

2. [使用 pprof 和火焰圖調試 golang 應用](https://cizixs.com/2017/09/11/profiling-golang-program/)
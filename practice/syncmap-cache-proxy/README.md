# Description

使用 sync.WaitGroup + sync.Map + closure 原理實現防止緩存穿透的方式

但還是無法完全解決大量併發請求湧入時，對資料庫請求不只一次的問題，只能說這個想法很妙

# References

- [Golang Taiwan Gathering #70](https://youtu.be/gwZhQiHcJlQ?t=2946)
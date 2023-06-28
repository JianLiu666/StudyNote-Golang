## Project Layout

```
StudyNode-Golang
 ├─ practice                                         # 實作指定情境
 │   ├─ game-management                                 # 遊戲模組 & 有限狀態機 (FSM)
 │   ├─ generics                                        # 泛型
 │   ├─ http2-server                                    # 支援 HTTP2 protocol 的 Server
 │   ├─ json-converter                                  # JSON 格式解析
 │   ├─ message-broker                                  # 跨 goroutines 的 message broker
 │   ├─ rethinking-classical-concurrency-patterns       # go doc sync.Cond 的 PDF 實作練習
 │   ├─ simple-cache                                    # Linked List & Hash Table
 │   ├─ syncmap-cache-proxy                             # WaitGroup + Sync Map 實現防止緩存穿透
 │   └─ simple-cache                                    # WebSocket packge benchmark
 ├─ standard-library                                 # 學習標準函式庫
 │   ├─ benchmark                                       # 壓力測試 & pprof
 │   ├─ context                                         # Context 調用關係
 │   ├─ interface                                       # Interface 調用關係
 │   ├─ pointer                                         # 物件、函式與指標調用關係
 │   ├─ select                                          # Channel Select 使用方式
 │   ├─ sizeof                                          # 資料結構的記憶體保存空間
 │   └─ wire                                            # DI Code generation tool
 ├─ third-party                                      # 學習第三方套件
 │   ├─ arangodb                                        # ArangoDB Client Driver
 │   ├─ fx                                              # Uber fx: Dependency Injection
 │   ├─ jwt                                             # JWT
 │   ├─ nats                                            # Messaging System: NATS.io
 │   └─ redis-stream                                    # Messaging System: Redis Stream
 └─ README.md        
```
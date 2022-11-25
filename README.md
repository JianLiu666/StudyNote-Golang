## Project Layout

```
StudyNode-Golang
 ├─ practice/                # 實作指定情境
 │   ├─ json-converter          # JSON 格式解析
 │   ├─ game-management         # 遊戲模組 & 有限狀態機 (FSM)
 │   └─ simple-cache            # Linked List & Hash Table
 ├─ standard-library/        # 學習標準函式庫
 │   ├─ benchmark/              # 壓力測試 & pprof
 │   ├─ context/                # Context 調用關係
 │   ├─ interface/              # Interface 調用關係
 │   ├─ pointer/                # 物件、函式與指標調用關係
 │   └─ sizeof                  # 資料結構的記憶體保存空間
 ├─ third-party/             # 學習第三方套件
 │   ├─ arangodb/               # ArangoDB Client Driver
 │   ├─ fx/                     # Uber fx: Dependency Injection
 │   ├─ jwt/                    # JWT
 │   └─ nats/                   # Messaging System: NATS.io
 └─ README.md        
```
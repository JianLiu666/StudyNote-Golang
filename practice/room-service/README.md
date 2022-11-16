# 實作筆記

簡單紀錄一下目前在工作上遇到的問題  
主要是想解決未來持續增加遊戲時, 開發者能只關心最上層業務邏輯實作即可

### 資料結構

```
baseModule
 ├─ services             # 第三方套件(redisClient、nats...)
 ├─ settings             # 底層共用設定檔
 └─ gameMap              # 遊戲管理列表
     ├─ gameName            # 遊戲名稱
     ├─ settings            # 遊戲共用設定檔
     └─ themeMap            # 遊戲大廳管理列表
         ├─ settings           # 遊戲大廳共用設定檔
         └─ roomMap            # 遊戲房間管理列表
             └─ IRoom             # 遊戲房間介面
                 └─ fsm              # 有限狀態機
```

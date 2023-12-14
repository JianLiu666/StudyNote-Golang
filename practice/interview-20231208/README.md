# Trading System

- [Trading System](#trading-system)
  - [Goal](#goal)
    - [業務場景](#業務場景)
    - [撮合機制設計](#撮合機制設計)
      - [收斂問題](#收斂問題)
      - [資料結構](#資料結構)
      - [撮合流程](#撮合流程)
      - [可用性保證](#可用性保證)
  - [API 設計](#api-設計)
    - [RESTful APIs](#restful-apis)
      - [Pending Order](#pending-order)
        - [Endpoint](#endpoint)
        - [Request](#request)
        - [Response](#response)
      - [Get Order by filters](#get-order-by-filters)
        - [Endpoint](#endpoint-1)
        - [Request](#request-1)
        - [Response](#response-1)
      - [Get Transaction Log by filters](#get-transaction-log-by-filters)
        - [Endpoint](#endpoint-2)
        - [Request](#request-2)
        - [Response](#response-2)
  - [Project Layout](#project-layout)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Quick install instructions](#quick-install-instructions)
    - [Usage](#usage)
  - [TODOs](#todos)

---

## Goal

### 業務場景

設計一個交易系統，主要功能如下：

1. 提供 Client 創建新訂單(e.g. buy/sell, market/limit)
2. 自動對訂單進行即時撮合，撮合機制以 FOK(filled-or-killed) 為主
3. 用戶可以根據訂單編號查閱訂單狀態 (e.g. 已成交, 已建立, etc.)

### 撮合機制設計

#### 收斂問題

- 搓合機制以 FOK 為主，只有當數量完全符合時才能撮合成功
- Limit Price: 只有當買方價格 >= 賣方價格時，才能嘗試撮合
  - 排序以 price 為主，當 price 相同時以 timestamp 為主 (ASC)
- Market Price: 直接以賣方最低點/買方最高點的價格嘗試撮合
  - 排序以 timestamp 為主 (ASC)
- Market Orders 的撮合順序優先於 Limit Orders
- 一旦撮合成功後會立即生成一筆交易紀錄

#### 資料結構

用兩個 Priority Queue (Heap) 維護買賣雙方的訂單排序

 - 買方: Max Heap
 - 賣方: Min Heap

一旦有新訂單建立時，立即比較兩個 PQ 的 top value 是否相同，如果相同表示撮合成功

 - 加入新訂單的時間複雜度 ~= `O(logn)`
 - 嘗試撮合的時間複雜度 ~= `O(1)`
   - 只需比較 top value

#### 撮合流程
 1. 接收訂單: 
    - 收到 limit order 後，加入到對應的 PQ
    - 收到 market order 時直接進行**撮合檢查**
 2. 撮合檢查: 
    - 有 market order 時直接與另一方的 PQ 撮合
    - 沒有 market order 時檢查 PQ 的 top value 是否滿足撮合條件: buyer price >= seller price
 3. 執行撮合:
    - 條件成立時扣除成交量(買賣方中較小的 quantity)，並且建立一筆交易紀錄
    - 清算完之後，清除 quantity 歸零的 orders

例外狀況:

 - 當市場上完全沒有任何的 limit orders，只剩下 market orders 時，會忽略所有的 market orders

#### 可用性保證

- 一但訂單狀態發生改變時，必須保證該更新同時寫入 Database
- 當系統啟動時，從 Database 讀取並復原所有未成交的訂單，繼續等待撮合

---

## API 設計

### RESTful APIs

#### Pending Order

發出掛單請求，等待系統根據規則進行撮合

##### Endpoint

```
[POST] /api/v1/orders
```

##### Request

- JSON Schema
```json
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
        "userId": {
            "type": "integer",
            "description": "用戶唯一識別碼"
        },
        "roleType": {
            "type": "integer",
            "description": "掛單角色 (0:買方, 1:賣方)"
        },
        "orderType": {
            "type": "integer",
            "description": "交易單類型 (0:市價單, 1:限價單)"
        },
        "durationType": {
            "type": "integer",
            "description": "交易單期限 (0:ROD, 1:IOC, 2:FOK)"
        },
        "price": {
            "type": "integer",
            "description": "交易單價格"
        },
        "quantity": {
            "type": "integer",
            "description": "交易單數量"
        }
    },
    "required": ["userId", "roleType", "orderType", "durationType", "price", "quantity"]
}

```

- Example
```json
{
    "userId": 1,
    "roleType": 0,
    "orderType": 0,
    "durationType": 0,
    "price": 100,
    "quantity": 100
}
```

##### Response

- Status Code
```
200: 請求成功
400: 參數錯誤
```

#### Get Order by filters

根據篩選條件查詢交易單

##### Endpoint

```
[GET] /api/v1/orders
```

##### Request

- Params
```
- `userId`         (int, optional): 指定用戶名稱
- `status`         (int, optional): 指定交易單狀態(0:掛單中, 1:已取消, 2:已完成)
- `startTimestamp` (timestamp, optional): 開始時間
- `endTimestamp`   (timestamp, optional): 結束時間
- `limit`          (int, optional): 查詢筆數上限
```

- Example
```
curl -X GET http://localhost:6600/api/v1/orders: 取回所有交易單
curl -X GET http://localhost:6600/api/v1/orders?userId=1: 取回用戶識別碼為 1 的所有交易單
```

##### Response

- Status Code
```
200: 請求成功
400: 參數錯誤
500: 系統錯誤
```

- JSON Schema
```json
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "array",
    "items": {
        "type": "object",
        "properties": {
            "id": {
                "type": "integer",
                "description": "交易單唯一識別碼"
            },
            "userId": {
                "type": "integer",
                "description": "用戶唯一識別碼"
            },
            "roleType": {
                "type": "integer",
                "description": "掛單角色 (0:買方, 1:賣方)"
            },
            "orderType": {
                "type": "integer",
                "description": "交易單類型 (0:市價單, 1:限價單)"
            },
            "durationType": {
                "type": "integer",
                "description": "交易單期限 (0:ROD, 1:IOC, 2:FOK)"
            },
            "price": {
                "type": "integer",
                "description": "交易單價格"
            },
            "quantity": {
                "type": "integer",
                "description": "交易單數量"
            },
            "status": {
                "type": "integer",
                "description": "交易單狀態 (0:掛單中, 1:已取消, 2:已完成)"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "交易單時間戳"
            }
        },
        "required": ["id", "userId", "roleType", "orderType", "durationType", "price", "quantity", "status", "timestamp"]
    }
}
****
```

- Example
```json
[
  {
    "id": 1,
    "userId": 1,
    "roleType": 0,
    "orderType": 1,
    "durationType": 0,
    "price": 100,
    "quantity": 100,
    "status": 0,
    "timestamp": "2023-12-14T23:38:39+08:00"
  },
  {
    "id": 2,
    "userId": 1,
    "roleType": 0,
    "orderType": 1,
    "durationType": 0,
    "price": 100,
    "quantity": 100,
    "status": 0,
    "timestamp": "2023-12-14T23:38:42+08:00"
  }
]
```

#### Get Transaction Log by filters

根據篩選條件查詢交易紀錄

##### Endpoint

```
[GET] /api/v1/transactions
```

##### Request

- Params
```
- `buyerOrderId`   (int, optional): 指定買方交易單唯一識別碼
- `sellerOrderId`  (int, optional): 指定賣方交易單唯一識別碼
- `startTimestamp` (timestamp, optional): 開始時間
- `endTimestamp`   (timestamp, optional): 結束時間
- `limit`          (int, optional): 查詢筆數上限
```

- Example
```
curl -X GET http://localhost:6600/api/v1/transactions: 取回所有交易單
curl -X GET http://localhost:6600/api/v1/transactions?buyerOrderId=1: 取回買方交易單唯一識別碼為 1 的所有交易單
```

##### Response

- Status Code
```
200: 請求成功
400: 參數錯誤
500: 系統錯誤
```

- JSON Schema
```json
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "array",
    "items": {
        "type": "object",
        "properties": {
            "id": {
                "type": "integer",
                "description": "成交紀錄唯一識別碼"
            },
            "buyerOrderId": {
                "type": "integer",
                "description": "買方交易單唯一識別碼"
            },
            "sellerOrderId": {
                "type": "integer",
                "description": "賣方交易單唯一識別碼"
            },
            "price": {
                "type": "integer",
                "description": "成交價格"
            },
            "quantity": {
                "type": "integer",
                "description": "成交數量"
            },
            "timestamp": {
                "type": "string",
                "format": "date-time",
                "description": "成交時間戳"
            }
        },
        "required": ["id", "buyerOrderId", "sellerOrderId", "price", "quantity", "timestamp"]
    }
}

```

- Example
```json
[
  {
    "id": 1,
    "buyerOrderId": 4,
    "sellerOrderId": 3,
    "price": 100,
    "quantity": 100,
    "timstamp": "2023-12-14T23:45:02+08:00"
  },
  {
    "id": 2,
    "buyerOrderId": 5,
    "sellerOrderId": 6,
    "price": 100,
    "quantity": 100,
    "timstamp": "2023-12-15T00:00:59+08:00"
  }
]
```

---

## Project Layout

- 參考 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

```
Trading System
 ├─ api/                # OpenAPI
 │   ├─ router/            # router group
 |   |   ├─ v1/               # v1 版本 APIs
 |   |   └─ router.go         # router common interface
 │   └─ server.go          # gin web framework
 ├─ cmd/                # 主要應用程式進入點
 ├─ config/             # 組態設定檔
 ├─ deployment/         # 部署設定檔
 ├─ model/              # Data schema
 ├─ pkg/                # 模組化函式庫
 │   ├─ accessor/          # 基礎建設管理模組 (e.g. config, network, storage, etc.)
 │   ├─ config/            # 組態設定模組 (viper)
 │   ├─ e/                 # 專案內部使用的狀態碼、型別定義
 │   ├─ rdb/               # Relational Database 模組
 │   └─ trading/           # 交易撮合模組
 ├─ dockerfile          #
 ├─ go.mod              #
 ├─ go.sum              #
 ├─ main.go             #
 ├─ makefile            #
 └─ README.md           #
```

---

## Getting Started

### Prerequisites

- Go
- Docker

### Quick install instructions

```shell
make init
```

### Usage

快速啟動 Trading 系統

```shell
make up
```

關閉系統

```shell
make down
```

單元測試

```shell
make test
```

---

## TODOs

- Features
  - [ ] 淘汰已經過期的 ROD orders
- Availability
  - [ ] Server Crash Recovery Mechanism
    - 從 MySQL 中恢復掛單中的 orders
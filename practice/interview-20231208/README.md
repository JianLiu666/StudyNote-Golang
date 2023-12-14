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
```

##### Response

- Status Code
```
200: 請求成功
400: 參數錯誤
```

- JSON Schema
```json
```

- Example
```json
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
```

##### Response

- Status Code
```
200: 請求成功
400: 參數錯誤
```

- JSON Schema
```json
```

- Example
```json
```

---

## Project Layout

- 參考 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

```
Trading System
 ├─ main.go             #
 └─ README.md           #
```

---

## Getting Started

TODO

---

## TODOs

- Features
  - [ ] 淘汰已經過期的 ROD orders
- Availability
  - [ ] Server Crash Recovery Mechanism
    - 從 MySQL 中恢復掛單中的 orders
# Trading System

- [Trading System](#trading-system)
  - [Goal](#goal)
    - [業務場景](#業務場景)
    - [撮合機制設計](#撮合機制設計)
      - [收斂問題](#收斂問題)
      - [資料結構](#資料結構)
      - [可用性保證](#可用性保證)
  - [API 設計](#api-設計)
    - [RESTful APIs](#restful-apis)
  - [Project Layout](#project-layout)
  - [Getting Started](#getting-started)

---

## Goal

### 業務場景

設計一個交易系統，主要功能如下：

1. 提供 Client 創建新訂單(Buy or Sell)
2. 根據已經生成的訂單進行撮合，用戶可以根據訂單編號查閱訂單狀態 (e.g. 已成交, 已建立, etc.)

### 撮合機制設計

#### 收斂問題

- 只有當買方與賣方的價格一致時，才能成交
- 訂單排序以 price 為主，當 price 相同時以 timestamp 為主 (ASC)
- 一旦撮合成功後會立即生成一筆交易紀錄

#### 資料結構

用兩個 Priority Queue (Heap) 維護買賣雙方的訂單排序

 - 買方: Max Heap
 - 賣方: Min Heap

一旦有新訂單建立時，立即比較兩個 PQ 的 top value 是否相同，如果相同表示撮合成功

 - 加入新訂單的時間複雜度 ~= `O(logn)`
 - 嘗試撮合的時間複雜度 ~= `O(1)`
   - 只需比較 top value

#### 可用性保證

- 一但訂單狀態發生改變時，必須保證該更新同時寫入 Database
- 當系統啟動時，從 Database 讀取並復原所有未成交的訂單，繼續等待撮合

---

## API 設計

### RESTful APIs

TODO

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
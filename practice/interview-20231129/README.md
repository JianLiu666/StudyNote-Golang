# Tinder Matching System

- [Tinder Matching System](#tinder-matching-system)
  - [Goal](#goal)
    - [業務場景](#業務場景)
    - [配對機制](#配對機制)
      - [規則設計](#規則設計)
      - [定義資料集](#定義資料集)
      - [設計配對模組](#設計配對模組)
      - [API 設計](#api-設計)
  - [Project Layout](#project-layout)
  - [Getting Started](#getting-started)

---

## Goal

### 業務場景

設計一個 Tinder-like 配對系統，主要功能如下:

1. AddSinglePersonAndMatch: 加入新用戶且根據**配對規則**進行配對與更新用戶狀態
2. RemoveSinglePerson: 主動移除用戶
3. QuerySinglePerson: 根據查詢條件返回符合條件的用戶

### 配對機制

#### 規則設計

- 男生只能配對到比自己還矮的女生；換句話說，女生只能配對到比自己還要高的男生
- 每個用戶都有自己的約會次數上限，一旦次數歸零就會立即從系統中移除

#### 定義資料集

- 根據規則，必須同時符合兩個條件才能視為成功配對：
  1. 必須是異性
  2. 男生身高必須大於女生
- 根據上述條件，對資料進行初步定義：
  1. 按照性別分類資料集
  2. 按照身高進行排序(asending order)
- 承上，一旦找到第一筆合法的資料就能持續索引下一筆資料，直到滿足任一條件：
  1. 用戶的約會次數已達上限
  2. 所有符合條件的用戶皆已配對

#### 設計配對模組

- 需求收斂：
  - 以姓名作為 Unique Key，避免出現重複用戶
  - 資料集在所有操作(CRUD) 中都需要保持有序
- 選擇資料結構：
  - Array：
    - 找到第一筆滿足匹配條件的時間複雜度是 `O(logn)` → Binary Search
    - 只要新增、刪除一筆資料，重新排序的時間複雜度是 `O(nlogn)`
  - Linked List：
    - 查詢的複雜度是 `O(n)`
    - 資料排序的複雜度可以降到 `O(n)` → 優於 Array
  - Hash Table：
    - 資料成映射關係(key-value pair)，不適合範圍查詢
  - Red-Black Tree：
    - 基於 Binary Search Tree
      - 有效避免當資料寫入是以有序寫入時，BST 退化成 Linked List 的問題
    - 查詢、插入、刪除的時間複雜度皆約為 `O(logn)`
- 基於上述考量，在實作中會選擇基於 Red-Black Tree 為主的 TreeMap 作為主要的資料結構保存資料
  - Key: 以 `{身高}-{姓名}` 作為 Composite Key 處理相同身高的重複問題
  - Value: 用戶資訊

#### API 設計

- TODO

---

## Project Layout

TODO

---

## Getting Started

TODO

---
# Key-Value 文章列表系統

- [Key-Value 文章列表系統](#key-value-文章列表系統)
  - [Goal](#goal)
    - [業務場景](#業務場景)
    - [改變文章資料的保存方式](#改變文章資料的保存方式)
      - [Why?](#why)
      - [What?](#what)
      - [How?](#how)
    - [資料庫選擇](#資料庫選擇)
      - [Redis](#redis)
  - [Features](#features)
    - [RESTful APIs](#restful-apis)
      - [GetHead](#gethead)
        - [Endpoint](#endpoint)
        - [Response](#response)
      - [GetPage](#getpage)
        - [Endpoint](#endpoint-1)
        - [Response](#response-1)
      - [Set](#set)
        - [Endpoint](#endpoint-2)
        - [Request Body](#request-body)
  - [Project Layout](#project-layout)
  - [References](#references)

---

## Goal

### 業務場景

- 每個用戶會有屬於自己的**個人化推薦**文章列表
- 負責處理推薦算法的團隊每小時會更新一次文章列表，因此系統的流量可能會很大
- 每個列表的內容只需要保留一天，如何有效率的清除過期的列表資料


### 改變文章資料的保存方式

#### Why?

 - 每當用戶發表一篇新文章時，系統會將文章寫進資料庫保存
 - 根據用戶偏好/訂閱看板不同，系統需要推薦對應的文章給不同用戶
   - 每當用戶繼續往下捲動頁面，系統就跟著繼續推送文章(i.e. Lazy loading)
   - 標準的分頁查詢情境，根據條件向資料庫查詢符合條件且第 `[i, j]` 筆的文章
   - e.g. 最新文章, 熱門文章, 精選推薦, etc.
 - 隨著資料規模/業務需求逐漸成長，RDB 的維護/查詢成本也會逐漸增加
   - 需要維護越來越多的 materialized views

#### What?

 - 文章間的關聯性是根據外部提供的規則作用在不同的 metadata 上所賦予；換句話說，可以定時主動建立文章間的關聯性
   - 將計算分頁的負載從取得文章列表的 SQL 中移出
   - 讓業務流程更加明確，**建立文章關聯性** 與 **維護分頁數量** 可以由不同的小組專門負責處理

#### How?

- System Design
```
              Set page   +--------------------+         +----------+
    Producer ----------> | Linked List Server | ------> | Database |
                         +--------------------+         +----------+
                                   |
                                   | Get head or specified page
                                   |
                                   v
                                  User
```

- Data structure
```
    e.g. 熱門文章
      
      topic key
         |
         v
      +------+            +------+       +------+                 +------+
      | page | ---------> | page | ----> | page | ----> ... ----> | page |
      +------+            +------+       +------+                 +------+
       ├─ articles
       └─ next page key
```

- 運用 Linked List 的概念維護文章列表的分頁，當用戶第一次訪問主題(e.g. 熱門文章, 個人推薦, etc.) 時根據這個主題對應的 `topic key` 取得 head page
- 隨著用戶持續瀏覽相同主題，就可以根據 next page key 指向拿到下一個 page 的文章識別碼(UUID)
- `Producer` 只需要專注在產生對應主題的 page content，由 `Linked List Server` 維護同一主題的資料結構與排序

### 資料庫選擇

#### Redis

 - Key-value store
 - 支援 JSON 格式
 - 可以透過 key expiration 維護文章列表的有效時間

---

## Features

### RESTful APIs

#### GetHead

取得指定主題的 head page uuid

##### Endpoint

```
[Get] /api/v1/head?listKey={string}
```

##### Response 

 - JSON Schema
```json
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
        "nextPageKey": {
            "type": "string"
        }
    },
    "required": ["nextPageKey"]
}
```
 - Example
```json
{
    "nextPageKey": "abcd"
}
```

#### GetPage

取得指定 page 的文章內容

##### Endpoint

```
[Get] /api/v1/page?pageKey={string}
```

##### Response

- JSON Schema
```json
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
        "articles": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "integer"
                    }
                },
                "required": ["id"]
            }
        },
        "nextPageKey": {
            "type": "string"
        }
    },
    "required": ["articles", "nextPageKey"]
}

```

- Example
```json
{
    "articles": [
        {
            "id": 123334
        },
        {
            "id": 123335
        }
    ],
    "nextPageKey": "efgh"
}
```

#### Set

對指定主題更新 head page

##### Endpoint

```
[Post] /api/v1/head
```

##### Request Body

- JSON Schema
```json
{
    "$schema": "https://json-schema.org/draft/2020-12/schema",
    "type": "object",
    "properties": {
        "listKey": {
            "type": "string"
        },
        "articles": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "id": {
                        "type": "integer"
                    }
                },
                "required": ["id"]
            }
        }
    },
    "required": ["listKey", "articles"]
}
```

- Example
```json
{
    "listKey": "hot",
    "articles": [
        {
            "id": 123334
        },
        {
            "id": 123335
        }
    ]
}
```

---

## Project Layout

```
TODO
```

---

## References
 - [Dcard Backend Team 如何讓工程師能更專注在列表排序與組合的演算法？](https://medium.com/dcardlab/dcard-backend-team-%E5%A6%82%E4%BD%95%E8%AE%93%E5%B7%A5%E7%A8%8B%E5%B8%AB%E8%83%BD%E6%9B%B4%E5%B0%88%E6%B3%A8%E5%9C%A8%E5%88%97%E8%A1%A8%E6%8E%92%E5%BA%8F%E8%88%87%E7%B5%84%E5%90%88%E7%9A%84%E6%BC%94%E7%AE%97%E6%B3%95-de07f45295f6)
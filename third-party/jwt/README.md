# 理解筆記

JWT 的組合可以看成是三個 JSON object 並且用 `.` 來做區隔, 各自編碼後組成一個 JWT 字串

- Header
- Payload
- Signature

## Header

必要欄位：

- alg：JWT 的主要加密演算法, 若是未加密則設置為 `none` (e.g. HMAC、SHA256、RSA ...)

非必要欄位：

- typ：JWT 本身媒體類型, 少數情境可能出現 JWT 與其他 JOSE header 混合使用的情況, 正常情境下即為 JWT
- cty：內容類型, 大多數情境下 JWT 會攜帶特定聲明與任意數據作為 payload 的一部分, 此時不得設置內容類型聲明, 因此 cty 聲明極少出現在 header 中

範例如下, 最後用 Base64 重新編碼：

```json
{
    "alg": "HS256",
    "typ": "JWT",
}
```

## Payload

這裡放的是聲明 (Claim) 的內容, 在定義上有三種聲明：

- Registered claims  
  標準公認的訊息, **建議**但不強迫遵守：
  - iss(Issuer)：發證者
  - sub(Subject)：主題
  - aud(Audience)：目標收件人
  - iat(Issued At)：簽發時間
  - exp(Expiration Time)：有效期限, 必須大於 JWT 簽發時間 (Unix Time)
  - nbf(Not Before)：多久之後 JWT 才開始正式生效 (Unix Time)
  - jti(JWT Id)：JWT UUID

- Private claims  
   自定義欄位, 配合實務需求增加 (e.g. UserAccount、UserName ...)

- Public claims  
   允許用戶在 [IANA JSON Web Token](https://www.iana.org/assignments/jwt/jwt.xhtml) 上註冊聲明, 實務上基本不會使用

通常所有使用者感興趣的資訊都會放在 payload 內, 如同紀錄在 session 內的用戶資訊, 範例如下：

```json
{
    "sub": "1234567890",
    "account": "JianLiu666@github.com",
    "role": "admin",
}
```

## Signature

由三大部份組成：

- base64UrlEncode(header)
- base64UrlEncode(payload)
- secret

<br/>

# Reference

1. [JWT(JSON Web Token)-原理介紹](https://kennychen-blog.herokuapp.com/2019/12/14/JWT-JSON-Web-Tokens-%E5%8E%9F%E7%90%86%E4%BB%8B%E7%B4%B9/)
2. [什麼是 JWT ?](https://5xruby.tw/posts/what-is-jwt/)
3. [淺談 JWT 的安全性與使用情境](https://medium.com/mr-efacani-teatime/%E6%B7%BA%E8%AB%87jwt%E7%9A%84%E5%AE%89%E5%85%A8%E6%80%A7%E8%88%87%E9%81%A9%E7%94%A8%E6%83%85%E5%A2%83-301b5491b60e)
4. [JWT 官網](https://jwt.io/)  
   官網有提供線上工具可以解析目前的 JWT 字串內容
5. [Golang-JWT 示範](https://medium.com/%E4%BC%81%E9%B5%9D%E4%B9%9F%E6%87%82%E7%A8%8B%E5%BC%8F%E8%A8%AD%E8%A8%88/golang-json-web-tokens-jwt-olang-json-web-tokens-jwt-%E7%A4%BA%E7%AF%84-225b377e0f79)
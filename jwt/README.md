# 理解筆記

JWT 的組合可以看成是三個 JSON object 並且用 `.` 來做區隔, 各自編碼後組成一個 JWT 字串
 - Header
 - Payload
 - Signature

## Header
 
1. alg
紀錄加密演算法
e.g. HMAC、SHA256、RSA ...

2. typ

<br/>

# Reference

1. [JWT(JSON Web Token)-原理介紹](https://kennychen-blog.herokuapp.com/2019/12/14/JWT-JSON-Web-Tokens-%E5%8E%9F%E7%90%86%E4%BB%8B%E7%B4%B9/)
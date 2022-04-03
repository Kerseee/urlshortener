# UrlShortener

UrlShortener is a simple http web application that provide url shortening and redirection. It is built from scratch in Go.

## Prerequisites
- [PostgreSQL](https://www.postgresql.org/)
- [GNU make](https://www.gnu.org/software/make/)
- [migrate](https://github.com/golang-migrate/migrate)

## Getting Start
1. Create a PostgreSQL database for UrlShortener.
2. Copy the database DSN and replace the value of $URLSHORTENER_DB_DSN in "./config/.envrc"
```
export URLSHORTENER_DB_DSN='postgres://urlshortener:password@localhost/urlshortener?sslmode=disable' 
```
3. Run database migrations
```
make db/migrations/up
```
4. Build UrlShortener
```
make build
```
5. Execute UrlShortener. The default server address is at localhost:8080.
```
./bin/urlshortner
```

## Usage
To shorten a url, please use <strong>curl</strong> to send a <strong>JSON</strong>-encoded request with <strong>POST</strong> method to the endpoint "http://{hostname:port}/api/v1/urls", and provide exactly <strong>"url"</strong> and <strong>"expireAt"</strong> these two fields.
```
curl -i -X POST -H 'Content-Type:application/json' -d '{"url":"http://github.com","expireAt":"2025-12-22T12:00:00Z"}' http://localhost:8080/api/v1/urls
```

If the request is valid, the client will receive a response like:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Sun, 03 Apr 2022 07:56:50 GMT
Content-Length: 68

{
	"id": "BQAwqbKa",
	"shortUrl": "http://localhost:8080/BQAwqbKa"
}
```

To redirect to the origin url, just GET the shortUrl:
```
curl -i -X GET http://localhost:8080/BQAwqbKa
```
Then the client will get the response:
```
HTTP/1.1 303 See Other
Content-Type: text/html; charset=utf-8
Location: http://github.com
Date: Sun, 03 Apr 2022 08:01:39 GMT
Content-Length: 44

<a href="http://github.com">See Other</a>.
```

### Request constraints
A valid request must contain a valid http or https url and an after-now expire time in valid JSON format. It should meet these constraints:
- Has exactly one "url" key and its value is a single string having prefix "http://" or "https://".
- Has exactly one "expireAt" key and it has a single JSON-formatted time value.
- Value of "expireAt" should not be before now.
- A request body should contain exactly one JSON object.

If one of the constraint is violated, the client will recieve a response like:
```
HTTP/1.1 400 Bad Request
Content-Type: application/json
Date: Sun, 03 Apr 2022 08:30:56 GMT
Content-Length: 69

{
	"error": [
		"invalid url",
		"expired time should after now"
	]
}
```

## Configuration
To configure the UrlShortener, please use following tags when starting the application:
|Tag|Usage|Type|Default|Notes|
|---|---|---|---|---|
|-h|Print flags||||
|-addr|Server address|string|localhost:8080||
|-db|Database DSN|string|$URLSHORTENER_DB_DSN||
|-db-max-idle-conns|Database maximum idle connections|int|25||
|-db-max-idle-time |Database maximum idle time|int|15|unit: minute|
|-db-max-open-conns|Database maximum open connections|int|25||
|-db-query-timeout |Database maximum query time|int|3|unit: second|
|-len-short-url|Length of shortened URL|int|8|should be greater than 4 and less than 17|
|-max-len-reshort-url|Maximum length of shortened URL for reshortening URL in case of short URL conflicts|int|12|should be greater than len-short-url and less than 44|

___

# 專案思路

### 第三方套件選擇
在開始專案之前，曾考慮是否使用第三方套件實作許多部分，例如後端框架 [Gin](https://github.com/gin-gonic/gin)、 servemux 套件 [httprouter](https://github.com/julienschmidt/httprouter)、資料庫連接 [Gorm](https://gorm.io/)、[測試 testify](https://github.com/stretchr/testify)、logging [zap](https://github.com/uber-go/zap) ，確實基於效能及維護難易程度考量，使用上述框架或第三方套件可能為較佳的選擇。

然而此專案作為一學習專案，且開發目的為應徵實習生職缺，恰好是一個更加深入學習 Go standard library 並藉此打好基礎的機會，個人認為若能先從 build from scratch 出發，未來需要使用框架或第三方套件時學習效率將會較為提升。

基於上述原因，此專案除了因 standard library 中 database/sql 必須使用第三方 driver 而使用了 [pq](github.com/lib/pq)，其餘皆無使用任何第三方套件。

## 專案架構
專案架構參考 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)、[Let's go](https://lets-go.alexedwards.net/) 以及 [Let's go further](https://lets-go-further.alexedwards.net/)。

## API 設計

### Request 前處理
於 pdf 中所述，一個正確的 request 如下
```
curl -X POST -H "Content-Type:application/json" http://localhost/api/v1/urls -d '{ 
    "url": "<original_url>", 
    "expireAt": "2021-02-08T09:20:41Z" 
}' 
```
然而因為 client 傳來的 request 有各種可能性，在縮網址時會擋掉以下幾種 request 並回傳 400 bad request：
- Method 不是 Post
- "url" 或 "expireAt" 的值為空值
- "url" 或 "expireAt" 的值型態錯誤
- "url" 的值不是 "http://" 或 "https://" 開頭
- "expireAt" 的時間早於現在
- JSON 裡面有其他的 field
- Request body 裡面 JSON 物件超過一個
- 其他 JSON syntax error

### Database schema
在此專案資料庫中只有一個 table: urls，欄位如下：

| Attribute | Type | Constraints |
| -------- | -------- | -------- |
| id     | bigserial     | primary key |
| url     | text     | not null |
| short_url     | text     | unique, not null |
|expire_at|time with time zone| not null|

考量 redirect 效能，在 short_url 上加了 unique constraint，並且加入 index (b-tree)。

### 如何縮網址
在構思如何縮網址時，考量了以下幾點：
1. 短網址應盡可能 unique，如此在 redirect 的 api 要從資料庫 IO 時比較有效率
2. 短網址應使人眼容易閱讀
3. 短網址長度應盡量統一
4. 短網址應夠短


如果要達成 1. 可以使用 hash function，例如 HD5, SHA-1, SHA-2 等，並使用 base64 或 base62 encoding 達成 2. 和 3.，最後取前 8 碼達成 4.。

然而在取前 8 碼時，原本 hash 完 unique 的特性會消失，導致兩個不同的原網址轉成同樣短網址的可能性上升，雖然機率仍然很低，但當流量越大時越可能發生。

考慮上述幾個因素，最後決定以 SHA-256 先 hash 原網址、使用 base64 encoding 並取前 8 碼，若遇到 conflict 則多取 1 碼，再次嘗試 insert，持續逐次增加短網址長度直到等於 -max-len-reshort-url 時若仍 conflict（雖然機率極低）才結束。

除了上述因素，若 request 中的原網址曾經縮過的話，在資料庫中也會有紀錄，然而因 short_url 設為 unique，且上述的縮網址方法並沒有 random 機制，這樣的 request 也會 conflict。為了解決這樣的 conflict，最後的縮網址 API 流程設計如下：

1. 驗證 request
2. 用 sha-256 及 base64 encoding 縮網址並取前 n 碼（預設為 8）
3. 試圖 insert 進 DB，若發生 conflict 則進入 step 4，若無則進入 step 7
4. 從資料庫裡取出同一個短網址的該筆資料，確認該筆資料的原網址與此 request 是否相同。若不相同則進入 step 5，相同則進入 step 6
5. 對於 hash 及 encoding 完後的短網址多取 1 碼，再次 insert 確認是否 conflict。若無則進入 step 7，若有則再次執行 step 5 直到短網址長度等於 -max-len-reshort-url
6. 確認 request 的 expire time 是否比資料庫中的同網址 expire time 晚，若是則更新資料庫中該筆資料的 expire time，若否則直接進入 step 7
7. 將短網址回應給 client

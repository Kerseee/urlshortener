# UrlShortener

UrlShortener is a simple http web application that provide url shortening and redirection. It is built from scratch in Go and PostgreSQL.

## Prerequisites
- [PostgreSQL](https://www.postgresql.org/)
- [GNU make](https://www.gnu.org/software/make/)
- [migrate](https://github.com/golang-migrate/migrate)

## Getting Start
1. Create a PostgreSQL database for UrlShortener.
2. Copy the database DSN and replace the $URLSHORTENER_DB_DSN in "./config/.envrc"
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
- Value of "expireAt" should not before now.
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

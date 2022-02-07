# Shortener

`Shortener` is a service to shorten long URLs to generated short URLs and handle the redirection.

> Study work, not yet perfect

## Config

**config/config.json**

```json
{
  "server": {
    "host": "127.0.0.1",
    "port": ":8080"
  },
  "redis": {
    "addr": "localhost:6379",
    "password": "",
    "db": 0
  },
  "exp": 30
}
```

## API

### Redirect

- Address: `/[token]`
- Example: `http://localhost:8080/MxQ9ycozO`
- Method: `GET`
- Description: `Redirect generated url to origin url`

### Generate

- Address: `/api/v1/generate`
- Example: `http://localhost:8080/api/v1/generate`
- Method: `POST`
- Description: `Long url generate short url with expiration`
- Request body: JSON with following parameters

| Parameters | Type   | Necessity | Description                                                          | Default    |
|------------|--------|-----------|----------------------------------------------------------------------|------------|
| url        | string | must      | long to shorten                                                      | none       |
| exp        | int    | suggest   | expiration in days<br>-1 means no expiration time| config.exp |

### Expire

- Address: `/api/v1/expire`
- Example: `http://localhost:8080/api/v1/expire`
- Method: `POST`
- Description: `Update token expiration time`
- Request body: JSON with following parameters

| Parameters | Type   | Necessity | Description        | Default |
|------------|--------|-----------|--------------------|---------|
| token      | string | must      | token for short URL| none    |
| exp        | int    | must      | expiration in days<br>-1 means no expiration time | none    |

## Usage

```go
package main

import "github.com/SukiEva/shortener"

func main() {
	if s, err := shortener.New(); err == nil {
		s.Serve()
	}
}
```

## Credits

- [Gin](https://github.com/gin-gonic/gin): A HTTP web framework written in Go (Golang)
- [Go-redis](https://github.com/go-redis/redis): Type-safe Redis client for Golang
- [Token](https://github.com/marksalpeter/token/): A simple base62 encoded token library for go, ideal for short url services

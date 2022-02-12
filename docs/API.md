# API

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

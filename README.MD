# Exchange Rate Api

## Test
You can run tests using
-   `go test ./... -v`

## Run
You should create config.yaml inside config folder, for example:
```yaml
Server:
  Port: 3000
ExchangeRateService:
  Urls:
    - https://run.mocky.io/v3/f05a9177-5719-47b3-94d9-60ab75e538a7
    - https://run.mocky.io/v3/7dcc2ac7-85f0-4032-a421-f8d947e20824
    - https://run.mocky.io/v3/377b087e-56bf-4a3f-a2f8-d662b4782705
ExchangeRateController:
  CacheDurationMinute: 1

```
After that you can run server using:
-   `go run .`

## Endpoints

## Get Exchange Rate
Get best exchange rate
-   method: `GET`
-   path: `/currency/{name}`
-   response:
    ```js
    {
    "code": integer, 
    "err": string,
    "data": ExchangeRateResponse{
        currency: string,
        value: float64,
    },
    }
    ```
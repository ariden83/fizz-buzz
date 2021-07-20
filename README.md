# fizz-buzz
Write a simple fizz-buzz REST server

## Installation for local

1. Install [go-swagger](https://github.com/go-swagger/go-swagger)
2. Launch API with `make local`

## Command

- `make local`          - launch servers
- `make local-test`     - launch tests
- `make local-bench`    - launch benches

## Example

#####  API

    http://127.0.0.1:8080/fizz-buzz?limit=100
    http://127.0.0.1:8080/fizz-buzz?limit=100&nbOne=3&nbTwo=5&strOne=fizz&strTwo=buzz

#####  Metrics, healthz, readiness

    http://127.0.0.1:8081/metrics
    http://127.0.0.1:8081/healthz
    http://127.0.0.1:8081/readiness

#####  Interface swagger
 
    http://127.0.0.1:8082

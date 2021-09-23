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

## Bench

#### Without cache

```
BenchmarkApi/Test_GET_/fizz-buzz?limit=100
BenchmarkApi/Test_GET_/fizz-buzz?limit=100-5         	    9044	    120842 ns/op	   19529 B/op	     449 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000-5        	    2516	    447132 ns/op	   40639 B/op	    2789 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000-5       	   21370	     50870 ns/op	   17416 B/op	     244 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000-5      	   23270	     48004 ns/op	   17238 B/op	     244 allocs/op
```

#### With cache

```
BenchmarkApi/Test_GET_/fizz-buzz?limit=100
BenchmarkApi/Test_GET_/fizz-buzz?limit=100-5         	   18226	     87179 ns/op	   17413 B/op	     246 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000-5        	   20994	     54943 ns/op	   17413 B/op	     246 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000-5       	   22284	     54849 ns/op	   17772 B/op	     245 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000-5      	   26648	     48535 ns/op	   17705 B/op	     245 allocs/op
```



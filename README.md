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
BenchmarkApi/Test_GET_/fizz-buzz?limit=100-5         	   10988	     97395 ns/op	   19539 B/op	     450 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000-5        	    2636	    459463 ns/op	   40706 B/op	    2790 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000-5       	     304	   4071531 ns/op	  271419 B/op	   26191 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000-5      	   19369	     52648 ns/op	   17726 B/op	     245 allocs/op
```

#### With cache

```
BenchmarkApi/Test_GET_/fizz-buzz?limit=100
BenchmarkApi/Test_GET_/fizz-buzz?limit=100-5         	   17506	     83726 ns/op	   17414 B/op	     246 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000
BenchmarkApi/Test_GET_/fizz-buzz?limit=1000-5        	   17914	     65442 ns/op	   17413 B/op	     246 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000
BenchmarkApi/Test_GET_/fizz-buzz?limit=10000-5       	   23389	     55954 ns/op	   17412 B/op	     246 allocs/op
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000
BenchmarkApi/Test_GET_/fizz-buzz?limit=100000-5      	   21924	     51687 ns/op	   17803 B/op	     245 allocs/op
```



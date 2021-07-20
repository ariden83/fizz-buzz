package benches

import (
	"net/http"
	"testing"
)

func (tts *Tests) GetFizzBuzz10000Bench(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	r, _ := http.NewRequest("GET", "/fizz-buzz?limit=10000&strOne=fizz&strTwo=buzz&nbOne=3&nbTwo=5", nil)
	benchRequest(b, tts.HTTPEndRouter, r)
}

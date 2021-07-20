package tests

import (
	"ariden/fizz-buzz/internal/catcher"
	"ariden/fizz-buzz/internal/endpoint"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type htmlResponse string

const (
	xRequestIDForTests string = "62661e3f-53f9-4bdb-b29f-4c6372dfd4a4"
	validPath          string = "/fizz-buzz"
)

var getFizzBuzzTests = []Scenario{
	{
		`Should be ok with default content-type "text/plain"`,
		validPath,
		200,
		`
		{
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100"
		}`,
		func(t *testing.T, args ...interface{}) {
			fizzBuzzResp := args[0]
			const waitingResp string = "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100"
			if fizzBuzzResp != waitingResp {
				t.Fatal("Bad response, have '", fizzBuzzResp, "' and we want '", waitingResp, "'")
			}
		},
		func(t *testing.T, header http.Header) {
			if xRequestId := header.Get("X-Request-Id"); xRequestId != xRequestIDForTests {
				t.Fatal("Fail to get Header X-Request-Id, have '", xRequestId, "' and we want '", xRequestIDForTests, '"')
			}
		},
	},
	{
		`Should be ok with default content-type "text/plain" AND "x-Request-Id" header in lowercase`,
		validPath,
		200,
		`
		{
			"x-request-id": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {
			if xRequestId := header.Get("X-Request-Id"); xRequestId != xRequestIDForTests {
				t.Fatal("Fail to get Header X-Request-Id, have '", xRequestId, "' and we want '", xRequestIDForTests, '"')
			}
		},
	},
	{
		`JSON: Should fail if limit is smaller than 1`,
		validPath,
		412,
		`
		{
			"x-request-id": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "0"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok if limit equal to 1`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"x-request-id": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "1"
		}`,
		func(t *testing.T, args ...interface{}) {
			fizzBuzzResp := (*args[0].(*endpoint.JsonResp))
			const waitingResp string = "1"
			if fizzBuzzResp.Txt != waitingResp {
				t.Fatal("Bad response, have '", fizzBuzzResp.Txt, "' and we want '", waitingResp, "'")
			}
		},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with default content-type "application/json"`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		``,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {
			if xRequestId := header.Get("X-Request-Id"); xRequestId != xRequestIDForTests {
				t.Fatal("Fail to get Header X-Request-Id, have '", xRequestId, "' and we want '", xRequestIDForTests, '"')
			}
		},
	},
	{
		`JSON: Should be ok with "limit" parameter`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100"
		}`,
		func(t *testing.T, args ...interface{}) {
			fizzBuzzResp := (*args[0].(*endpoint.JsonResp))
			const waitingResp string = "1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,62,63,64,65,66,67,68,69,70,71,72,73,74,75,76,77,78,79,80,81,82,83,84,85,86,87,88,89,90,91,92,93,94,95,96,97,98,99,100"
			if fizzBuzzResp.Txt != waitingResp {
				t.Fatal("Bad response, have '", fizzBuzzResp.Txt, "' and we want '", waitingResp, "'")
			}
		},
		func(t *testing.T, header http.Header) {
			if xRequestId := header.Get("X-Request-Id"); xRequestId != xRequestIDForTests {
				t.Fatal("Fail to get Header X-Request-Id, have '", xRequestId, "' and we want '", xRequestIDForTests, '"')
			}
		},
	},
	{
		`JSON: Should be ok with "limit", "nbOne", "strOne" parameters`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"nbOne": "2",
			"strOne": "fizz"
		}`,
		func(t *testing.T, args ...interface{}) {
			fizzBuzzResp := (*args[0].(*endpoint.JsonResp))
			const waitingResp string = "1,fizz,3,fizz,5,fizz,7,fizz,9,fizz,11,fizz,13,fizz,15,fizz,17,fizz,19,fizz,21,fizz,23,fizz,25,fizz,27,fizz,29,fizz,31,fizz,33,fizz,35,fizz,37,fizz,39,fizz,41,fizz,43,fizz,45,fizz,47,fizz,49,fizz,51,fizz,53,fizz,55,fizz,57,fizz,59,fizz,61,fizz,63,fizz,65,fizz,67,fizz,69,fizz,71,fizz,73,fizz,75,fizz,77,fizz,79,fizz,81,fizz,83,fizz,85,fizz,87,fizz,89,fizz,91,fizz,93,fizz,95,fizz,97,fizz,99,fizz"
			if fizzBuzzResp.Txt != waitingResp {
				t.Fatal("Bad response, have '", fizzBuzzResp.Txt, "' and we want '", waitingResp, "'")
			}
		},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with "limit", "nbTwo", "strTwo" parameters`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"nbTwo": "3",
			"strTwo": "buzz"
		}`,
		func(t *testing.T, args ...interface{}) {
			fizzBuzzResp := (*args[0].(*endpoint.JsonResp))
			const waitingResp string = "1,2,buzz,4,5,buzz,7,8,buzz,10,11,buzz,13,14,buzz,16,17,buzz,19,20,buzz,22,23,buzz,25,26,buzz,28,29,buzz,31,32,buzz,34,35,buzz,37,38,buzz,40,41,buzz,43,44,buzz,46,47,buzz,49,50,buzz,52,53,buzz,55,56,buzz,58,59,buzz,61,62,buzz,64,65,buzz,67,68,buzz,70,71,buzz,73,74,buzz,76,77,buzz,79,80,buzz,82,83,buzz,85,86,buzz,88,89,buzz,91,92,buzz,94,95,buzz,97,98,buzz,100"
			if fizzBuzzResp.Txt != waitingResp {
				t.Fatal("Bad response, have '", fizzBuzzResp.Txt, "' and we want '", waitingResp, "'")
			}
		},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with "limit", "nbOne", "nbTwo", "strOne", "strTwo" parameters`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"nbOne": "1",
			"nbTwo": "2",
			"strOne": "fizz",
			"strTwo": "buzz"
		}`,
		func(t *testing.T, args ...interface{}) {
			fizzBuzzResp := (*args[0].(*endpoint.JsonResp))
			const waitingResp string = "fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz,fizz,fizzbuzz"
			if fizzBuzzResp.Txt != waitingResp {
				t.Fatal("Bad response, have '", fizzBuzzResp.Txt, "' and we want '", waitingResp, "'")
			}
		},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with "limit" equal to "1000000000"`,
		validPath,
		412,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "1000000000",
			"nbOne": "1",
			"nbTwo": "2",
			"strOne": "fizz",
			"strTwo": "buzz"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should fail with "nbOne" and without "strOne" parameters`,
		validPath,
		412,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"nbOne": "1"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should fail with "nbTwo" and without "strTwo" parameters`,
		validPath,
		412,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"nbTwo": "1"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should fail with only "nbOne", "nbTwo" and without "str" parameters`,
		validPath,
		412,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"nbTwo": "1",
			"nbOne": "2"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with "strTwo" and without "nbTwo" parameters`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"strTwo": "toto"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with "strOne" and without "nbOne" parameters`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"strOne": "toto"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok with "strOne", "two", and without "nbOne", "nbTwo" parameters`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `",
			"X-Request-ID": "` + xRequestIDForTests + `"
		}
		`,
		`{
			"limit": "100",
			"strOne": "toto",
			"strOne": "two"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {},
	},
	{
		`JSON: Should be ok without "X-Request-ID" header`,
		validPath,
		200,
		`
		{
			"Content-Type": "` + endpoint.ContentTypeJSON + `"
		}
		`,
		`{
			"limit": "100"
		}`,
		func(t *testing.T, args ...interface{}) {},
		func(t *testing.T, header http.Header) {
			if xRequestId := header.Get("X-Request-Id"); xRequestId == xRequestIDForTests {
				t.Fatal("Fail to get Header X-Request-Id, have '", xRequestId, "' and we don't want '", xRequestIDForTests, '"')
			}
		},
	},
}

type toto struct {
	Plus string
}

func (tts *Tests) GetFizzBuzzTest(t *testing.T) {
	for _, test := range getFizzBuzzTests {
		t.Run(test.description, func(t *testing.T) {
			client := &http.Client{}
			var err error
			URL, err := tts.getURL(test)
			if err != nil {
				t.Fatal("fail to get URL of unit test", err.Error())
			}

			r, err := http.NewRequest(http.MethodGet, URL, nil)
			if err != nil {
				t.Fatal("fail to GET ", err.Error())
			}

			if err := setHeaders(r, test); err != nil {
				t.Fatal("fail to set headers ", err.Error())
			}

			response, err := client.Do(r)
			if err != nil {
				t.Fatal(err.Error())
			}
			defer response.Body.Close()

			if test.expectedHeader != nil {
				test.expectedHeader(t, response.Header)
			}

			if response.StatusCode != test.statusCode {
				t.Fatal("wrong http status returned ", response.StatusCode, ", we want ", test.statusCode, URL)
			}

			if test.expectedBody != nil {

				buffer, err := ioutil.ReadAll(response.Body)
				if err != nil {
					t.Fatal("error with ioutil.ReadAll in CallPostTest")
				}

				if strings.Index(r.Header.Get("Content-Type"), endpoint.ContentTypeJSON) != -1 {

					appRetrieved := &endpoint.JsonResp{}
					json.Unmarshal(buffer, appRetrieved)
					catcher.Block{
						Try: func() {
							test.expectedBody(t, appRetrieved)
						},
						Catch: func(e catcher.Exception) {
							if e != nil {
								t.Logf("Caught %v\n", e)
							}
							var err error
							json.Unmarshal(buffer, err)
							if err != nil {
								t.Fatalf("have error %s", fmt.Sprintf("%#v", err))
							}
						},
					}.Do()

				} else {
					test.expectedBody(t, string(buffer))
				}

			}
		})
	}
}

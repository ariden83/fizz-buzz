package endpoint

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Resp string

type JsonResp struct {
	Txt string `json:"txt"`
}

// getFizzBuzzResp screen response
//
// swagger:response getFizzBuzzResp
// nolint
type getFizzBuzzResp struct {
	// Content-Length
	// in: header
	ContentLength string `json:"Content-Length"`
	// Content-Type
	// in: header
	ContentType string `json:"Content-Type"`
	// X-Request-Id
	// in: header
	XRequestID string `json:"X-Request-Id"`
	// corps of Response
	// in: body
	Body JsonResp `json:"body"`
}

// getFizzBuzzReq Params for method GET
//
// swagger:parameters getFizzBuzzReq
// nolint
type getFizzBuzzReq struct {
	// Content-Type
	// in: header
	ContentType string `json:"Content-Type"`
	// X-Request-Id
	// in: header
	XRequestID string `json:"X-Request-Id"`
	getFizzBuzzParams
}

type getFizzBuzzParams struct {
	// Number one
	// in: query
	NBOne int `json:"nbOne"`
	// Number two
	// in: query
	NBTwo int `json:"nbTwo"`
	// limit
	// in: query
	Limit int `json:"limit"`
	// String One
	// in: query
	StrOne string `json:"strOne"`
	// String two
	// in: query
	StrTwo string `json:"strTwo"`
	isJSON bool
}

// getFizzBuzz swagger:route GET /fizz-buzz fizzbuzz getFizzBuzzReq
//
// Get fizzBuzz filters by 5 parameters
//
//     Consumes:
//     - application/json
//     - text/html
//
//     Produces:
//     - application/json
//     - text/html
//
//     Schemes: http, https
//
// Responses:
//    default: genericError
//        200: getFizzBuzzResp
//        401: genericError
//        404: genericError
//        412: genericError
//        500: genericError
func (m *Endpoint) GetFizzBuzz(w http.ResponseWriter, r *http.Request, _ map[string]string) {
	params := getFizzBuzzParams{}
	if err := m.checkRequest(&params, r); err != nil {
		m.fail(http.StatusPreconditionFailed, err, w, r)
		return
	}

	defer m.IncMetrics(params)
	if m.cache != nil {
		cacheKey := fmt.Sprintf("%v", params)
		resp := m.cache.Fetch(cacheKey, func() []byte {
			ch := make(chan string, 1)
			go m.convert(ch, p)
			return m.formatEntireStringResp(ch, p.isJSON)
		})
		
		if _, err := w.Write(resp); err != nil {
			m.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
			m.fail(http.StatusInternalServerError, err, w, r)
		}
	}
	
	ch := make(chan string, 1)
	go m.convert(ch, params)
	m.formatResp(w, r, params, ch)
}


func (m *Endpoint) generateResponse (w http.ResponseWriter, r *http.Request, p getFizzBuzzParams) {
	ch := make(chan string, 1)
	go m.convert(ch, params)
	m.formatResp(w, r, params, ch)
}

func (m *Endpoint) IncMetrics(p getFizzBuzzParams) {
	m.metrics.ApiParamsCounter.WithLabelValues(strconv.Itoa(p.Limit), strconv.Itoa(p.NBOne), strconv.Itoa(p.NBTwo), p.StrOne, p.StrTwo).Inc()
}

func (m *Endpoint) checkRequest(p *getFizzBuzzParams, r *http.Request) error {
	var (
		err error
		q   url.Values = r.URL.Query()
	)

	if q.Get("nbOne") != "" {
		p.NBOne, err = strconv.Atoi(q.Get("nbOne"))
		if err != nil {
			return fmt.Errorf("invalid integer for p nbOne %s", q.Get("nbOne"))
		}
		if p.NBOne > m.conf.Parameters.MaxNb {
			return fmt.Errorf("maximum size exceeded for p NBOne %s, max %d", q.Get("nbOne"), m.conf.Parameters.MaxNb)
		}
		if p.NBOne == 0 {
			return fmt.Errorf("NBOne peter must be greater than zero")
		}
	}

	if q.Get("nbTwo") != "" {
		p.NBTwo, err = strconv.Atoi(q.Get("nbTwo"))
		if err != nil {
			return fmt.Errorf("invalid integer for p NBTwo %s", q.Get("nbTwo"))
		}
		if p.NBTwo > m.conf.Parameters.MaxNb {
			return fmt.Errorf("maximum size exceeded for p NBTwo %s, max %d", q.Get("nbTwo"), m.conf.Parameters.MaxNb)
		}
		if p.NBTwo == 0 {
			return fmt.Errorf("NBTwo peter must be greater than zero")
		}
	}

	if q.Get("limit") != "" {
		p.Limit, err = strconv.Atoi(q.Get("limit"))
		if err != nil {
			return fmt.Errorf("invalid integer for p limit %s", q.Get("limit"))
		}
		if p.Limit < 1 {
			return fmt.Errorf("limit peter must be greater zero")
		}
		if p.Limit > m.conf.Parameters.MaxLimit {
			return fmt.Errorf("maximum size exceeded for p limit %s, max %d", q.Get("limit"), m.conf.Parameters.MaxLimit)
		}
	}

	p.StrOne = q.Get("strOne")
	if len(p.StrOne) > m.conf.Parameters.MaxStrChar {
		return fmt.Errorf("maximum char exceeded %s, max %d", p.StrOne, m.conf.Parameters.MaxStrChar)
	}

	p.StrTwo = q.Get("strTwo")
	if len(p.StrTwo) > m.conf.Parameters.MaxStrChar {
		return fmt.Errorf("maximum char exceeded %s, max %d", p.StrTwo, m.conf.Parameters.MaxStrChar)
	}

	p.isJSON = strings.Index(r.Header.Get("Content-Type"), ContentTypeJSON) != -1

	return nil
}

// Returns a list of strings with numbers from 1 to limit,
// where:
// all multiples of int1 are replaced by str1,
// all multiples of int2 are replaced by str2,
// all multiples of int1 and int2 are replaced by str1str2.
func (Endpoint) convert(ch chan string, p getFizzBuzzParams) {
	defer close(ch)
	if p.Limit == 0 {
		ch <- ""
		return
	}
	for i := 1; i <= p.Limit; i++ {
		var (
			str        string
			isMultiple bool
		)
		if p.NBOne > 0 && i%p.NBOne == 0 {
			str = p.StrOne
			isMultiple = true
		}
		if p.NBTwo > 0 && i%p.NBTwo == 0 {
			str += p.StrTwo
			isMultiple = true
		}
		if !isMultiple {
			str = strconv.Itoa(i)
		}
		if i < p.Limit {
			str += ","
		}
		ch <- str
	}
	return
}

func (m *Endpoint) formatResp(w http.ResponseWriter, r *http.Request, p getFizzBuzzParams, ch chan string) {
	var finalJsonStr string = ""
	for {
		if result, ok := <-ch; ok {
			if p.isJSON {
				finalJsonStr += result
				continue
			}
			fmt.Fprint(w, result)
		} else {
			break
		}
	}

	if p.isJSON {
		resp := JsonResp{
			Txt: finalJsonStr,
		}
		if js, err := json.Marshal(resp); err != nil {
			m.log.Error("Fail to json.Marshal", zap.Error(err))
			m.fail(http.StatusInternalServerError, err, w, r)

		} else if _, err := w.Write(js); err != nil {
			m.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
			m.fail(http.StatusInternalServerError, err, w, r)
		}
	}

	return
}

func (m *Endpoint) formatEntireStringResp(ch chan string, isJson bool) []byte {
	var finalJsonStr string = ""
	for {
		if result, ok := <-ch; ok {
			finalJsonStr += result
		} else {
			break
		}
	}

	if isJson {
		resp := JsonResp{
			Txt: finalJsonStr,
		}
		js, _ := json.Marshal(resp)
		return js
	}

	return []byte(finalJsonStr)
}

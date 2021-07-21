package endpoint

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

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
	// required: true
	ContentLength string `json:"Content-Length"`
	// Content-Type
	// in: header
	// required: true
	ContentType string `json:"Content-Type"`
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// corps of Response
	// in: body
	Body JsonResp `json:"body"`
}

// getFizzBuzzReq ps for method GET
//
// swagger:parameters getFizzBuzzReq
// nolint
type getFizzBuzzReq struct {
	// Content-Type
	// in: header
	// required: true
	ContentType string `json:"Content-Type"`
	// X-Request-Id
	// in: header
	// required: true
	XRequestID string `json:"X-Request-Id"`
	// ps
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
}

type Resp string

// getFizzBuzz swagger:route GET /fizz-buzz fizz-buzz getFizzBUzzReq
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
	m.log.Info("Get fizz buzz")

	p := getFizzBuzzParams{}
	if err := m.checkRequest(&p, r.URL.Query()); err != nil {
		m.fail(http.StatusPreconditionFailed, err, w, r)
		return
	}

	defer m.IncMetrics(p)

	resp := convert(p)
	m.formatResp(w, r, resp)
}

func (m *Endpoint) IncMetrics(p getFizzBuzzParams) {
	m.metrics.ApiParamsCounter.WithLabelValues(strconv.Itoa(p.Limit), strconv.Itoa(p.NBOne), strconv.Itoa(p.NBTwo), p.StrOne, p.StrTwo).Inc()
}

func (m *Endpoint) checkRequest(p *getFizzBuzzParams, req url.Values) error {
	var err error
	if req.Get("nbOne") != "" {
		p.NBOne, err = strconv.Atoi(req.Get("nbOne"))
		if err != nil {
			return fmt.Errorf("invalid integer for p nbOne %s", req.Get("nbOne"))
		}
		if p.NBOne > m.conf.Parameters.MaxNb {
			return fmt.Errorf("maximum size exceeded for p NBOne %s, max %d", req.Get("nbOne"), m.conf.Parameters.MaxNb)
		}
		if p.NBOne == 0 {
			return fmt.Errorf("NBOne peter must be greater than zero")
		}
	}

	if req.Get("nbTwo") != "" {
		p.NBTwo, err = strconv.Atoi(req.Get("nbTwo"))
		if err != nil {
			return fmt.Errorf("invalid integer for p NBTwo %s", req.Get("nbTwo"))
		}
		if p.NBTwo > m.conf.Parameters.MaxNb {
			return fmt.Errorf("maximum size exceeded for p NBTwo %s, max %d", req.Get("nbTwo"), m.conf.Parameters.MaxNb)
		}
		if p.NBTwo == 0 {
			return fmt.Errorf("NBTwo peter must be greater than zero")
		}
	}

	if req.Get("limit") != "" {
		p.Limit, err = strconv.Atoi(req.Get("limit"))
		if err != nil {
			return fmt.Errorf("invalid integer for p limit %s", req.Get("limit"))
		}
		if p.Limit < 1 {
			return fmt.Errorf("limit peter must be greater zero")
		}
		if p.Limit > m.conf.Parameters.MaxLimit {
			return fmt.Errorf("maximum size exceeded for p limit %s, max %d", req.Get("limit"), m.conf.Parameters.MaxLimit)
		}
	}

	p.StrOne = req.Get("strOne")
	if len(p.StrOne) > m.conf.Parameters.MaxStrChar {
		return fmt.Errorf("maximum char exceeded %s, max %d", p.StrOne, m.conf.Parameters.MaxStrChar)
	}

	p.StrTwo = req.Get("strTwo")
	if len(p.StrTwo) > m.conf.Parameters.MaxStrChar {
		return fmt.Errorf("maximum char exceeded %s, max %d", p.StrTwo, m.conf.Parameters.MaxStrChar)
	}

	return nil
}

// Returns a list of strings with numbers from 1 to limit,
// where:
// all multiples of int1 are replaced by str1,
// all multiples of int2 are replaced by str2,
// all multiples of int1 and int2 are replaced by str1str2.
func convert(p getFizzBuzzParams) string {
	var finalStr string
	if p.Limit == 0 {
		return finalStr
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
		finalStr += str + ","
	}
	return finalStr[0 : len(finalStr)-1]
}

func (m *Endpoint) formatResp(w http.ResponseWriter, r *http.Request, response string) {

	if strings.Index(r.Header.Get("Content-Type"), ContentTypeJSON) != -1 {
		resp := JsonResp{
			Txt: response,
		}
		if js, err := json.Marshal(resp); err != nil {
			m.log.Error("Fail to json.Marshal", zap.Error(err))
			m.fail(http.StatusInternalServerError, err, w, r)

		} else if _, err := w.Write(js); err != nil {
			m.log.Error("Fail to Write response in http.ResponseWriter", zap.Error(err))
			m.fail(http.StatusInternalServerError, err, w, r)
		}
		return
	}
	fmt.Fprint(w, response)
	return
}

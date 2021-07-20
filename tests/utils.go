package tests

import (
	"ariden/fizz-buzz/config"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"
)

type Tests struct {
	Conf       *config.Config
	DefaultURL string
}

type test func(t *testing.T, args ...interface{})
type testHeader func(t *testing.T, header http.Header)

type Scenario struct {
	description    string
	route          string
	statusCode     int
	headers        string
	qs             string
	expectedBody   test
	expectedHeader testHeader
}

type ScenarioWithBody struct {
	description    string
	route          string
	statusCode     int
	body           string
	expectedBody   test
	expectedHeader testHeader
}

func pad(n int) string {
	if n < 10 {
		return fmt.Sprintf("0%d", n)
	}
	return fmt.Sprintf("%d", n)
}

// startFunctionnalTests start functionnal test when video respond
func (tts *Tests) StartFunctionnalTests(t *testing.T) {
	for {
		err := tts.callAPI(t)
		if err == nil {
			return
		}
		fmt.Println("Retrying to call API because of API fail with error :", err.Error())
		time.Sleep(1000 * time.Millisecond)
	}
}

// HealthCheckTest Test if healthz route respond ok
func (tts *Tests) HealthCheckTest(t *testing.T) {
	addr := fmt.Sprintf("http://%s:%d/healthz", tts.Conf.Host, tts.Conf.Metrics.Port)
	response, error := http.Get(addr)
	if error != nil {
		t.Fatal("An error occured when trying to check the healthz")
	}
	// THEN
	if response.StatusCode != http.StatusOK {
		t.Fatal("Wrong http status returned from healthcheck")
	}
}

// MetricsTest Test if metrics route respond ok
func (tts *Tests) MetricsTest(t *testing.T) {
	addr := fmt.Sprintf("http://%s:%d/metrics", tts.Conf.Host, tts.Conf.Metrics.Port)
	response, error := http.Get(addr)
	if error != nil {
		t.Fatal("An error occured when trying to check the metrics")
	}

	// THEN
	if response.StatusCode != http.StatusOK {
		t.Fatal("Wrong http status returned from metrics endpoint")
	}
}

type PingRes struct {
	Message string
}

func (tts *Tests) GetPingTest(t *testing.T) {
	t.Run("test ping HTTP method", func(t *testing.T) {
		client := &http.Client{}
		var err error
		var response *http.Response
		// generate URL
		URL, err := tts.getURL(Scenario{
			route: "/ping",
		})
		if err != nil {
			t.Fatal("fail to get URL of unit test", err.Error())
		}
		// Send POST request
		req, err := http.NewRequest(http.MethodGet, URL, nil)
		// test response
		if err != nil {
			t.Fatal("fail to GET ", err.Error())
		}
		response, err = client.Do(req)
		if err != nil {
			t.Fatal(err.Error())
		}
		defer response.Body.Close()

		buffer, err := ioutil.ReadAll(response.Body)
		if err != nil {
			t.Fatal("error with ioutil.ReadAll in CallPostTest")
		}

		row := &PingRes{}
		json.Unmarshal(buffer, row)
		if row.Message != "pong" {
			t.Fatal("Fail to call ping HTTP method, want pong response, have ", row)
		}
		if response.StatusCode != 200 {
			t.Fatal("wrong http status returned ", response.StatusCode, ", we want 200 with ping HTTP method")
		}
	})
}

// setHeaders set headers for test request
func setHeaders(req *http.Request, test Scenario) error {
	if test.headers != "" {
		var params = map[string]string{}
		if err := json.Unmarshal([]byte(test.headers), &params); err != nil {
			return err
		}
		for k := range params {
			req.Header.Add(k, params[k])
		}
	}
	return nil
}

// getURL formatted path + query in a valid url
// nolint
func (tts *Tests) getURL(test Scenario) (string, error) {
	var URL *url.URL

	URL, err := url.Parse(tts.DefaultURL)
	if err != nil {
		return "", err
	}
	// get each params in json string and add it to parameters
	URL.Path += test.route
	parameters := url.Values{}

	if test.qs != "" {
		var params = map[string]string{}
		if err = json.Unmarshal([]byte(test.qs), &params); err != nil {
			return "", err
		}
		for k := range params {
			parameters.Add(k, fmt.Sprintf("%s", params[k]))
		}
	}
	URL.RawQuery = parameters.Encode()
	return URL.String(), nil
}

// callAPI call video to know if it respond well
func (tts *Tests) callAPI(t *testing.T) error {
	// generate URL
	client := &http.Client{}

	URL, err := tts.getURL(Scenario{route: "fizz-buzz"})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	// test response
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode >= 500 {
		return errors.New(fmt.Sprintf("Video API not respond code 200, have code %d", response.StatusCode))
	}
	return nil
}

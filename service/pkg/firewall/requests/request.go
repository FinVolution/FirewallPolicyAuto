package requests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/FinVolution/FirewallPolicyAuto/service/pkg/logger"
)

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       []byte `json:"body"`
}

// HTTPClient Simple http client
type HTTPClient struct {
	client *http.Client
}

// HTTPClient url method and params struct
type RequestParams struct {
	URL       string
	Method    string
	BasicAuth struct {
		Username string
		Password string
	}
	QueryParams map[string]string
	Headers     map[string]string
	Body        any
}

// NewHTTPClient Create HTTPClient instance
func NewHTTPClient(insecureSkipVerify bool) *HTTPClient {
	return &HTTPClient{client: &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				PreferServerCipherSuites: true,
				InsecureSkipVerify:       insecureSkipVerify},
		},
	}}
}

// Serialize http request body to JSON
func (rp *RequestParams) marshalBody() ([]byte, error) {
	return json.Marshal(rp.Body)
}

func (rp *RequestParams) combineQueryURL() string {
	if len(rp.QueryParams) > 0 {
		query := url.Values{}
		for k, v := range rp.QueryParams {
			query[k] = []string{v}
		}
		return fmt.Sprintf("%s?%s", rp.URL, query.Encode())
	}
	return rp.URL
}

// Request http, support GET、POST etc. support Basic Authentication
func (hc *HTTPClient) Request(params *RequestParams) (*Response, error) {
	reqBody, err := params.marshalBody()
	if err != nil {
		return &Response{}, err
	}

	req, err := http.NewRequest(params.Method, params.combineQueryURL(), bytes.NewBuffer(reqBody))
	if err != nil {
		return &Response{}, err
	}

	// Set application json header
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// Set headers
	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	// Set up Basic Auth if username and password are given
	if params.BasicAuth.Username != "" && params.BasicAuth.Password != "" {
		req.SetBasicAuth(params.BasicAuth.Username, params.BasicAuth.Password)
	}

	// Send request
	resp, err := hc.client.Do(req)
	if err != nil {
		return &Response{}, err
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &Response{StatusCode: resp.StatusCode, Body: respBody}, err
	}
	return &Response{StatusCode: resp.StatusCode, Body: respBody}, nil
}

// RequestGetStatusCodeCheck  get request status code check
func RequestGetStatusCodeCheck(funcName string, body []byte, statusCode int) (bool, error) {

	if statusCode > http.StatusOK && statusCode < http.StatusMultipleChoices {
		logger.Infof("%s resp statusCode: %d", funcName, statusCode)
		return false, nil
	}

	if statusCode != http.StatusOK {
		err := fmt.Errorf("%s failed, statusCode: %d, resp body: %s", funcName, statusCode, string(body))
		logger.Error(err.Error())
		return false, err
	}
	return true, nil
}

// RequestPostStatusCodeCheck  post request status code check
func RequestPostStatusCodeCheck(funcName string, body []byte, statusCode, successCode int) error {
	if successCode == 0 {
		successCode = http.StatusCreated
	}
	if statusCode != successCode {
		err := fmt.Errorf("%s failed, statusCode: %d, resp body: %s", funcName, statusCode, string(body))
		logger.Error(err.Error())
		return err
	}
	return nil
}

package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RestAdapterTestSuite ...
type RestAdapterTestSuite struct {
	BaseURL string
	Client  *http.Client
	Headers map[string]string
}

// NewRestAdapterTestSuite ...
func NewRestAdapterTestSuite(baseURL string, client *http.Client, headers map[string]string) *RestAdapterTestSuite {
	h := make(map[string]string)
	h["Accept"] = "application/json;charset=UTF-8"
	h["Content-Type"] = "application/json"

	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			h[k] = v
		}
	}

	newRestAdapterTestSuite := &RestAdapterTestSuite{
		BaseURL: baseURL,
		Client:  client,
		Headers: h,
	}
	return newRestAdapterTestSuite
}

// ZPost ...
func (z *RestAdapterTestSuite) ZPost(url string, payload interface{}, requiredStatusCode int) (interface{}, error) {
	return z.httpSender(http.MethodPost, url, payload, requiredStatusCode)
}

// ZPut ...
func (z *RestAdapterTestSuite) ZPut(url string, payload interface{}, requiredStatusCode int) (interface{}, error) {
	return z.httpSender(http.MethodPut, url, payload, requiredStatusCode)
}

// ZGet ...
func (z *RestAdapterTestSuite) ZGet(url string, requiredStatusCode int) (interface{}, error) {
	return z.httpSender(http.MethodGet, url, nil, requiredStatusCode)
}

// ZDelete ...
func (z *RestAdapterTestSuite) ZDelete(url string, requiredStatusCode int) (interface{}, error) {
	return z.httpSender(http.MethodDelete, url, nil, requiredStatusCode)
}

// ZUpsert ...
func (z *RestAdapterTestSuite) ZUpsert(url string, payload interface{}, requiredStatusCode int) (interface{}, error) {
	return z.httpSender(http.MethodPost, url, payload, requiredStatusCode)
}

func (z *RestAdapterTestSuite) httpSender(httpMethod string, url string, payload interface{}, requiredStatusCode int) (interface{}, error) {
	var err error
	var b []byte
	b = nil

	if payload != nil {
		b, err = json.MarshalIndent(payload, "", "    ")
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(httpMethod, fmt.Sprintf("%s%s", z.BaseURL, url), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	response, err := z.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != requiredStatusCode {
		b, err := readResponseBody(response)
		if err == nil {
			var tmp interface{}
			err = json.Unmarshal(b, &tmp)
			if err != nil {
				return nil, err
			}
			return tmp, fmt.Errorf("expecting status code %d; got %d", requiredStatusCode, response.StatusCode)
		}
		return nil, fmt.Errorf("expecting status code %d; got %d", requiredStatusCode, response.StatusCode)
	}

	b, err = readResponseBody(response)
	if err != nil {
		return nil, err
	}

	var tmp interface{}
	err = json.Unmarshal(b, &tmp)
	if err != nil {
		return nil, err
	}
	return tmp, nil
}

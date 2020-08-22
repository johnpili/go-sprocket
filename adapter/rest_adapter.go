package adapter

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// RestAdapter ...
type RestAdapter struct {
	BaseURL string
	Client  *http.Client
	Headers map[string]string
}

// NewRestAdapter ...
func NewRestAdapter(baseURL string, client *http.Client, headers map[string]string) *RestAdapter {
	h := make(map[string]string)
	h["Accept"] = "application/json;charset=UTF-8"
	h["Content-Type"] = "application/json"

	if headers != nil && len(headers) > 0 {
		for k, v := range headers {
			h[k] = v
		}
	}

	newRestAdapter := &RestAdapter{
		BaseURL: baseURL,
		Client:  client,
		Headers: h,
	}
	return newRestAdapter
}

// ZPost ...
func (z *RestAdapter) ZPost(url string, payload interface{}) (int, []byte, error) {
	return z.httpSender(http.MethodPost, url, payload)
}

// ZPut ...
func (z *RestAdapter) ZPut(url string, payload interface{}) (int, []byte, error) {
	return z.httpSender(http.MethodPut, url, payload)
}

// ZGet ...
func (z *RestAdapter) ZGet(url string) (int, interface{}, error) {
	return z.httpSender(http.MethodGet, url, nil)
}

// ZDelete ...
func (z *RestAdapter) ZDelete(url string) (int, interface{}, error) {
	return z.httpSender(http.MethodDelete, url, nil)
}

// ZUpsert ...
func (z *RestAdapter) ZUpsert(url string, payload interface{}) (int, []byte, error) {
	return z.httpSender(http.MethodPost, url, payload)
}

// ZGetOne ...
func (z *RestAdapter) ZGetOne() (int, []byte, error) {
	return http.StatusNotImplemented, nil, errors.New("Not implemented")
}

// ZGetMany ...
func (z *RestAdapter) ZGetMany() (int, []byte, error) {
	return http.StatusNotImplemented, nil, errors.New("Not implemented")
}

// ZExecute ...
func (z *RestAdapter) ZExecute() (int, error) {
	return http.StatusNotImplemented, errors.New("Not implemented")
}

// ExtractResponsePayload ...
//func (z *RestAdapter) ExtractResponsePayload(payload interface{}) (*api.ResponsePayload, error) {
//}

func (z *RestAdapter) httpSender(httpMethod string, url string, payload interface{}) (int, []byte, error) {
	var err error
	var b []byte
	b = nil

	if payload != nil {
		b, err = json.MarshalIndent(payload, "", "    ")
		if err != nil {
			return http.StatusInternalServerError, nil, err
		}
	}

	req, err := http.NewRequest(httpMethod, fmt.Sprintf("%s%s", z.BaseURL, url), bytes.NewBuffer(b))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	response, err := z.Client.Do(req)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	/*if response.StatusCode < 200 || response.StatusCode > 299 {
		tmp, err := readResponseBody(response)
		if err == nil {
			return tmp, fmt.Errorf("None 200 HTTP response code; got %d", response.StatusCode)
		}
		return nil, fmt.Errorf("None 200 HTTP response code; got %d", response.StatusCode)
	}*/

	body, err := readResponseBody(response)
	return response.StatusCode, body, err
}

func readResponseBody(response *http.Response) ([]byte, error) {
	return bodyReader(response.Body)
}

func readRequestBody(request *http.Request) ([]byte, error) {
	return bodyReader(request.Body)
}

func bodyReader(body io.ReadCloser) ([]byte, error) {
	defer body.Close()
	return ioutil.ReadAll(body)

	/*rawBytes, err := ioutil.ReadAll(body)

	if err != nil {
		return nil, err
	}

	var buffer interface{}
	err = json.Unmarshal(rawBytes, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil*/
}

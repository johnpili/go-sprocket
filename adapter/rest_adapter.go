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

var ()

// RestAdapter ...
type RestAdapter struct {
	baseURL   string
	transport *http.Transport
	client    *http.Client
}

// NewRestAdapter ...
func NewRestAdapter(baseURL string, tr *http.Transport) *RestAdapter {
	newRestAdapter := &RestAdapter{
		baseURL:   baseURL,
		transport: tr}
	return newRestAdapter
}

// ZPost ...
func (z *RestAdapter) ZPost(url string, payload interface{}) (interface{}, error) {
	return z.pppSender(http.MethodPost, url, payload)
}

// ZPut ...
func (z *RestAdapter) ZPut(url string, payload interface{}) (interface{}, error) {
	return z.pppSender(http.MethodPut, url, payload)
}

// ZGet ...
func (z *RestAdapter) ZGet(url string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

// ZDelete ...
func (z *RestAdapter) ZDelete(url string) (interface{}, error) {
	return nil, errors.New("Not implemented")
}

// ZUpsert ...
func (z *RestAdapter) ZUpsert(url string, payload interface{}) (interface{}, error) {
	return z.pppSender(http.MethodPost, url, payload)
}

// ZGetOne ...
func (z *RestAdapter) ZGetOne() (interface{}, error) {
	return nil, errors.New("Not implemented")
}

// ZGetMany ...
func (z *RestAdapter) ZGetMany() (interface{}, error) {
	return nil, errors.New("Not implemented")
}

// ZExecute ...
func (z *RestAdapter) ZExecute() error {
	return errors.New("Not implemented")
}

func (z *RestAdapter) pppSender(httpMethod string, url string, payload interface{}) (interface{}, error) {
	b, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(httpMethod, fmt.Sprintf("%s%s", z.baseURL, url), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	response, err := z.client.Do(req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 200 && response.StatusCode <= 299 {
		return nil, fmt.Errorf("None 200 HTTP response code; got %s", response.StatusCode)
	}

	return z.readResponseBody(response)
}

func (z *RestAdapter) readResponseBody(response *http.Response) (interface{}, error) {
	return z.bodyReader(response.Body)
}

func (z *RestAdapter) readRequestBody(request *http.Request) (interface{}, error) {
	return z.bodyReader(request.Body)
}

func (z *RestAdapter) bodyReader(body io.ReadCloser) (interface{}, error) {
	rawBytes, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, err
	}

	var buffer interface{}
	err = json.Unmarshal(rawBytes, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

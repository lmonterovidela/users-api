package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type RestClient interface {
	Get(url string, headers map[string]string, body interface{}) (*http.Response, error)
}

type RestProvider struct {
	BaseURL string
	Client  *http.Client
}

func (r *RestProvider) Get(url string, headers map[string]string, body interface{}) (*http.Response, error) {
	buffer := bytes.Buffer{}
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("can't marshal body: %v, err:%v", body, err))
		}
		buffer = *bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(http.MethodGet, r.BaseURL+url, &buffer)

	if headers != nil {
		for key, value := range headers {
			req.Header.Add(key, value)
		}
	}

	if err != nil {
		return nil, err
	}
	return r.Client.Do(req)
}

func NewRestClient(baseUrl string, timeout time.Duration) RestClient {
	client := &http.Client{
		Timeout: timeout * time.Millisecond,
	}

	return &RestProvider{BaseURL: baseUrl, Client: client}
}

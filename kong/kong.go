package kong

import (
	"net/http"
	"encoding/json"
	"bytes"
	"io/ioutil"
	"errors"
)

type AuthHeader struct {
	Key string
	Value string
}

type Client struct {
	httpClient *http.Client
	authHeader *AuthHeader
	path string
}

func Make(httpClient *http.Client, header *AuthHeader, path string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: httpClient,
		authHeader: header,
		path: path,
	}
}

func (c Client) AddCertificate(cert, key, domain string) error {
	js, _ := json.Marshal(map[string]interface{}{
		"cert": cert,
		"key": key,
		"snis": []string{domain}})

	rq, err := http.NewRequest("POST", c.path + "/certificates/", bytes.NewReader(js))
	if err != nil {
		return err
	}
	rq.Header.Add("Content-Type", "application/json")
	if c.authHeader != nil {
		rq.Header.Add(c.authHeader.Key, c.authHeader.Value)
	}

	rsp, err := c.httpClient.Do(rq)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode >= 300 {
		bts, _ := ioutil.ReadAll(rsp.Body)
		return errors.New(string(bts))
	}

	return nil
}

type Certificate struct {
	ID string `json:"id"`
	Cert string `json:"cert"`
	Key string `json:"key"`
	SNIS []string `json:"snis"`
	CreatedAt int `json:"created_at"`
}

type CertListResponse struct {
	Total int `json:"total"`
	Data []Certificate `json:"data"`
}

func (c Client) GetCertificates() (*CertListResponse, error) {
	rq, err := http.NewRequest("GET", c.path + "/certificates/", nil)
	if err != nil {
		return nil, err
	}
	rq.Header.Add("Accept", "application/json")

	if c.authHeader != nil {
		rq.Header.Add(c.authHeader.Key, c.authHeader.Value)
	}

	rsp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode >= 300 {
		bts, _ := ioutil.ReadAll(rsp.Body)
		return nil, errors.New(string(bts))
	}

	var certRsp CertListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&certRsp); err != nil {
		return nil, err
	}

	return &certRsp, nil
}
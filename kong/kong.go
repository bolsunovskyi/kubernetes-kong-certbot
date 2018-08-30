package kong

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type AuthHeader struct {
	Key   string
	Value string
}

type Client struct {
	httpClient *http.Client
	authHeader *AuthHeader
	path       string
}

func Make(httpClient *http.Client, header *AuthHeader, path string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		httpClient: httpClient,
		authHeader: header,
		path:       path,
	}
}

func (c Client) AddCertificate(cert, key, domain string) error {
	js, _ := json.Marshal(map[string]interface{}{
		"cert": cert,
		"key":  key,
		"snis": []string{domain}})

	rq, _ := http.NewRequest("POST", c.path+"/certificates/", bytes.NewReader(js))
	_, err := c.sendRequest(rq)

	return err
}

type Certificate struct {
	ID        string   `json:"id"`
	Cert      string   `json:"cert"`
	Key       string   `json:"key"`
	SNIS      []string `json:"snis"`
	CreatedAt int      `json:"created_at"`
}

type CertListResponse struct {
	Total int           `json:"total"`
	Data  []Certificate `json:"data"`
}

func (c Client) GetCertificates() (*CertListResponse, error) {
	rq, _ := http.NewRequest("GET", c.path+"/certificates/", nil)
	rsp, err := c.sendRequest(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var certRsp CertListResponse
	if err := json.NewDecoder(rsp.Body).Decode(&certRsp); err != nil {
		return nil, err
	}

	return &certRsp, nil
}

func (c Client) sendRequest(rq *http.Request) (*http.Response, error) {
	rq.Header.Add("Accept", "application/json")
	rq.Header.Add("Content-Type", "application/json")

	if c.authHeader != nil {
		rq.Header.Add(c.authHeader.Key, c.authHeader.Value)
	}

	rsp, err := c.httpClient.Do(rq)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode >= 300 {
		bts, _ := ioutil.ReadAll(rsp.Body)
		return nil, errors.New(string(bts))
	}

	return rsp, nil
}

func (c Client) GetCertificate(host string) (*Certificate, error) {
	rq, _ := http.NewRequest("GET", c.path+"/certificates/"+host, nil)
	rsp, err := c.sendRequest(rq)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var cert Certificate
	err = json.NewDecoder(rsp.Body).Decode(&cert)
	return &cert, err
}

func (c Client) UpdateCertificate(host, cert, key string) error {
	js, _ := json.Marshal(map[string]interface{}{
		"cert": cert,
		"key":  key,
	})

	rq, _ := http.NewRequest("PATCH", c.path+"/certificates"+host, bytes.NewReader(js))
	_, err := c.sendRequest(rq)
	return err
}

func (c Client) UpdateOrCreateCertificate(host, cert, key string) error {
	js, _ := json.Marshal(map[string]interface{}{
		"cert": cert,
		"key":  key,
		"snis": []string{host},
	})

	rq, _ := http.NewRequest("PUT", c.path+"/certificates"+host, bytes.NewReader(js))
	_, err := c.sendRequest(rq)
	return err
}

func (c Client) DeleteCertificate(host string) error {
	rq, _ := http.NewRequest("DELETE", c.path+"/certificates"+host, nil)
	_, err := c.sendRequest(rq)
	return err
}

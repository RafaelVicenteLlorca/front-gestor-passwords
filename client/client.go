package client

import (
	"crypto/tls"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	HttpClient = New()
)

type HTTPClientCustom struct {
	Client          *http.Client
	BackendUri      string
	ContentTypeJSON string
	token           string
}

func initTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
}

func initServer() *HTTPClientCustom {
	return &HTTPClientCustom{
		Client:          &http.Client{Transport: initTransport()},
		BackendUri:      os.Getenv("BACKEND_URI"),
		ContentTypeJSON: "application/json",
		token:           "",
	}
}

func New() *HTTPClientCustom {
	return initServer()
}

func (c *HTTPClientCustom) SetToken(token string) {
	c.token = token
}

func (c *HTTPClientCustom) GetToken() string {
	return c.token
}

func (c *HTTPClientCustom) Do(req *http.Request) (*http.Response, error) {
	if c.token != "" {
		req.Header.Add("Authorization", "Bearer "+c.token)
	}
	req.Header.Add("Content-Type", c.ContentTypeJSON)
	return c.Client.Do(req)
}

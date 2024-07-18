package thunes

import (
	"net/http"

	"github.com/gustavo-bordin/thunes/config"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ThunesClient struct {
	username   string
	password   string
	httpClient HTTPClient
	hostUrl    string
}

func NewClient(cfg config.CliConfig) ThunesClient {
	return ThunesClient{
		username:   cfg.Thunes.Username,
		password:   cfg.Thunes.Password,
		hostUrl:    cfg.Thunes.HostUrl,
		httpClient: &http.Client{},
	}
}

func (c ThunesClient) Do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.username, c.password)
	return c.httpClient.Do(req)
}

package ngrok

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gustavo-bordin/thunes/config"
)

type ngrok struct {
	containerUrl string
}

type tunnel struct {
	PublicUrl string `json:"public_url"`
}

type ngrokResponse struct {
	Tunnels []tunnel `json:"tunnels"`
}

func NewNgrok(cfg config.CliConfig) ngrok {
	return ngrok{
		containerUrl: cfg.Ngrok.Url,
	}
}

func (n ngrok) GetNgrokUrl() (*string, error) {
	resp, err := http.Get(n.containerUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get ngrok URL")
	}

	var ngrokResponse ngrokResponse
	if err := json.NewDecoder(resp.Body).Decode(&ngrokResponse); err != nil {
		return nil, err
	}

	if len(ngrokResponse.Tunnels) == 0 {
		return nil, errors.New("no tunnels found")
	}

	firstTunnel := 0
	url := ngrokResponse.Tunnels[firstTunnel].PublicUrl
	return &url, nil
}

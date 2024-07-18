package thunes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	availablePayersEndpoint = "/v2/money-transfer/payers/"
)

func (c ThunesClient) GetAvailablePayers() ([]Payer, error) {
	url := fmt.Sprintf("%s%s", c.hostUrl, availablePayersEndpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var payers []Payer
	err = json.NewDecoder(res.Body).Decode(&payers)
	if err != nil {
		return nil, err
	}

	return payers, nil
}

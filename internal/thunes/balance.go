package thunes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	getBalancesEndpoint = "/v2/money-transfer/balances"
)

func (c ThunesClient) GetBalances() ([]Balance, error) {
	url := fmt.Sprintf("%s%s", c.hostUrl, getBalancesEndpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var balances []Balance
	err = json.NewDecoder(res.Body).Decode(&balances)
	if err != nil {
		return nil, err
	}

	return balances, nil
}

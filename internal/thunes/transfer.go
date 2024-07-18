package thunes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	createQuotationEndpoint    = "/v2/money-transfer/quotations"
	createTransactionEndpoint  = "/v2/money-transfer/quotations/%d/transactions"
	confirmTransactionEndpoint = "/v2/money-transfer/transactions/ext-%s/confirm"
)

func (c ThunesClient) CreateQuotation(q CreateQuotationDto) (*Quotation, error) {
	url := fmt.Sprintf("%s%s", c.hostUrl, createQuotationEndpoint)
	jsonPayload, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	bytesPayload := bytes.NewBuffer(jsonPayload)

	req, err := http.NewRequest(http.MethodPost, url, bytesPayload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var quotation Quotation
	err = json.NewDecoder(res.Body).Decode(&quotation)
	if err != nil {
		return nil, err
	}

	return &quotation, nil
}

func (c ThunesClient) CreateTransaction(
	payload CreateTransactionDto,
	quotationId int,
) (*Transaction, error) {
	endpoint := fmt.Sprintf(createTransactionEndpoint, quotationId)
	url := fmt.Sprintf("%s%s", c.hostUrl, endpoint)

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	bytesPayload := bytes.NewBuffer(jsonPayload)
	req, err := http.NewRequest(http.MethodPost, url, bytesPayload)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var transaction Transaction
	err = json.NewDecoder(res.Body).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (c ThunesClient) ConfirmTransaction(extId string) (*Transaction, error) {
	endpoint := fmt.Sprintf(confirmTransactionEndpoint, extId)
	url := fmt.Sprintf("%s%s", c.hostUrl, endpoint)

	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var transaction Transaction
	err = json.NewDecoder(res.Body).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

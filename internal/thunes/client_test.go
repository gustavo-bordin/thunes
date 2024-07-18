package thunes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestGetAvailablePayers(t *testing.T) {
	mockDoFunc := func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "/v2/money-transfer/payers/", req.URL.Path)
		assert.Equal(t, "Basic dGVzdF91c2VybmFtZTp0ZXN0X3Bhc3N3b3Jk", req.Header.Get("Authorization"))

		payers := []Payer{
			{ID: 1, Name: "Payer 1"},
			{ID: 2, Name: "Payer 2"},
		}
		resBody, _ := json.Marshal(payers)
		res := &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(bytes.NewBuffer(resBody)),
		}
		return res, nil
	}

	mockClient := &MockHTTPClient{
		DoFunc: mockDoFunc,
	}

	client := ThunesClient{
		username:   "test_username",
		password:   "test_password",
		httpClient: mockClient,
		hostUrl:    "http://mockserver.com",
	}

	payers, err := client.GetAvailablePayers()

	assert.NoError(t, err)

	totalPayersExpected := 2
	assert.Len(t, payers, totalPayersExpected)

	assert.Equal(t, int32(1), payers[0].ID)
	assert.Equal(t, "Payer 1", payers[0].Name)
	assert.Equal(t, int32(2), payers[1].ID)
	assert.Equal(t, "Payer 2", payers[1].Name)
}

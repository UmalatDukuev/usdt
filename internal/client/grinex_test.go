package client_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"usdt/internal/client"

	"github.com/stretchr/testify/assert"
)

func TestGetDepth_Success(t *testing.T) {
	mockResponse := `{
		"asks": [{"price": "100.5", "volume": "1", "amount": "1", "factor": "1", "type": "ask"}],
		"bids": [{"price": "99.5", "volume": "1", "amount": "1", "factor": "1", "type": "bid"}]
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	c := client.GrinexClient{URL: server.URL}
	rate, err := c.GetDepth()
	assert.NoError(t, err)
	assert.Equal(t, 100.5, rate.Ask)
	assert.Equal(t, 99.5, rate.Bid)
	assert.NotZero(t, rate.Timestamp)
}

func TestGetDepth_HTTPError(t *testing.T) {
	// Не существует такого сервера — вызовет ошибку http.Get
	c := client.GrinexClient{URL: "http://nonexistent.invalid"}
	rate, err := c.GetDepth()
	assert.Error(t, err)
	assert.Nil(t, rate)
}

func TestGetDepth_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("not json")); err != nil {
			log.Printf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	c := client.GrinexClient{URL: server.URL}
	rate, err := c.GetDepth()
	assert.Error(t, err)
	assert.Nil(t, rate)
}

func TestGetDepth_EmptyDepth(t *testing.T) {
	mockResponse := `{"asks": [], "bids": []}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("Failed to write mock response: %v", err)
		}
	}))
	defer server.Close()

	c := client.GrinexClient{URL: server.URL}
	rate, err := c.GetDepth()
	assert.Error(t, err)
	assert.Equal(t, "empty depth", err.Error())
	assert.Nil(t, rate)
}

func TestGetDepth_InvalidAskPrice(t *testing.T) {
	mockResponse := `{
		"asks": [{"price": "not_a_float", "volume": "1", "amount": "1", "factor": "1", "type": "ask"}],
		"bids": [{"price": "99.5", "volume": "1", "amount": "1", "factor": "1", "type": "bid"}]
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	c := client.GrinexClient{URL: server.URL}
	rate, err := c.GetDepth()
	assert.Error(t, err)
	assert.Nil(t, rate)
}

func TestGetDepth_InvalidBidPrice(t *testing.T) {
	mockResponse := `{
		"asks": [{"price": "100.5", "volume": "1", "amount": "1", "factor": "1", "type": "ask"}],
		"bids": [{"price": "not_a_float", "volume": "1", "amount": "1", "factor": "1", "type": "bid"}]
	}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte(mockResponse)); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	c := client.GrinexClient{URL: server.URL}
	rate, err := c.GetDepth()
	assert.Error(t, err)
	assert.Nil(t, rate)
}

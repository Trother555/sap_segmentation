package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sap_segmentation/internal/config"
	"sap_segmentation/internal/logger"

	"github.com/stretchr/testify/assert"
)

func TestFetchData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/test?p_limit=10&p_offset=0", r.URL.String())
		assert.Equal(t, "GET", r.Method)

		username, password, ok := r.BasicAuth()
		assert.True(t, ok)
		assert.Equal(t, "testuser", username)
		assert.Equal(t, "testpass", password)

		assert.Equal(t, "test-agent", r.Header.Get("User-Agent"))

		response := ErpResp{
			Segmentation: []*Segmentation{
				{
					DateFrom:     time.Now(),
					DateTo:       time.Now(),
					AddressSapId: "12345",
					AdrSegment:   "SegmentA",
					SegmentId:    1,
				},
			},
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	cfg := &config.Config{
		ConnURI:          server.URL + "/test",
		ConnTimeout:      1,
		ConnAuthLoginPwd: "testuser:testpass",
		ConnUserAgent:    "test-agent",

		LogPath: "test_log.txt",
	}

	logger := logger.New(cfg)

	client := New(cfg, logger)

	ctx := context.Background()
	result, err := client.FetchData(ctx, 0, 10)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, len(result.Segmentation))
	assert.Equal(t, "12345", result.Segmentation[0].AddressSapId)
	assert.Equal(t, "SegmentA", result.Segmentation[0].AdrSegment)
	assert.Equal(t, int64(1), result.Segmentation[0].SegmentId)
}

func TestFetchData_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := &config.Config{
		ConnURI:          server.URL + "/test",
		ConnTimeout:      1,
		ConnAuthLoginPwd: "testuser:testpass",
		ConnUserAgent:    "test-agent",

		LogPath: "test_log.txt",
	}

	logger := logger.New(cfg)

	client := New(cfg, logger)

	ctx := context.Background()
	result, err := client.FetchData(ctx, 0, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "api error status: 500 Internal Server Error")
}

func TestFetchData_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	cfg := &config.Config{
		ConnURI:          server.URL + "/test",
		ConnTimeout:      1,
		ConnAuthLoginPwd: "testuser:testpass",
		ConnUserAgent:    "test-agent",

		LogPath: "test_log.txt",
	}

	logger := logger.New(cfg)

	client := New(cfg, logger)

	ctx := context.Background()
	result, err := client.FetchData(ctx, 0, 10)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid character")
}

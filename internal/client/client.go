package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sap_segmentation/internal/config"
	"sap_segmentation/internal/logger"
	"strings"
	"time"
)

type Client interface {
	FetchData(ctx context.Context, offset int64, batchSize int64) (*ErpResp, error)
}

type ClientImpl struct {
	c        *http.Client
	logger   logger.Logger
	cfg      *config.Config
	username string
	password string
}

func New(cfg *config.Config, logger logger.Logger) Client {
	creds := strings.Split(cfg.ConnAuthLoginPwd, ":")

	return &ClientImpl{
		c: &http.Client{
			Timeout: time.Second * time.Duration(cfg.ConnTimeout),
		},
		logger:   logger,
		cfg:      cfg,
		username: creds[0],
		password: creds[1],
	}
}

type Segmentation struct {
	DateFrom     time.Time `json:"DATE_FROM"`
	DateTo       time.Time `json:"DATE_TO"`
	AddressSapId string    `json:"ADDRESS_SAP_ID"`
	AdrSegment   string    `json:"ADR_SEGMENT"`
	SegmentId    int64     `json:"SEGMENT_ID"`
}

type ErpResp struct {
	Segmentation []*Segmentation `json:"SEGMENTATION"`
}

func (c *ClientImpl) FetchData(ctx context.Context, offset int64, batchSize int64) (*ErpResp, error) {
	url := fmt.Sprintf("%s?p_limit=%d&p_offset=%d", c.cfg.ConnURI, batchSize, offset)
	c.logger.Info("fetching data from %s", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("User-Agent", c.cfg.ConnUserAgent)
	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error status: %s", res.Status)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fetchetResuld := &ErpResp{}
	err = json.Unmarshal(body, fetchetResuld)
	if err != nil {
		return nil, err
	}
	return fetchetResuld, nil
}

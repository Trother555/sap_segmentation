package importer

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"sap_segmentation/internal/client"
	"sap_segmentation/internal/config"
	"sap_segmentation/internal/logger"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockClient struct {
	mock.Mock
}

func (m *MockClient) FetchData(ctx context.Context, offset int64, batchSize int64) (*client.ErpResp, error) {
	args := m.Called(ctx, offset, batchSize)
	return args.Get(0).(*client.ErpResp), args.Error(1)
}

type MockDb struct {
	mock.Mock
}

func (m *MockDb) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	args := m.Called(ctx, query, arg)
	return nil, args.Error(1)
}

func (m *MockDb) GetDb() *sqlx.DB {
	return nil
}

func TestImporter_Run(t *testing.T) {
	mockClient := new(MockClient)
	mockDb := new(MockDb)

	cfg := &config.Config{
		ImportBatchSize:  10,
		ConnInterval:     1000,
		LogCleanupMaxAge: int64(time.Hour * 7),
		LogPath:          "test_log.txt",
	}

	logger := logger.New(cfg)

	importer := New(cfg, logger, mockClient, mockDb)

	mockClient.On("FetchData", mock.Anything, int64(1), int64(10)).Return(&client.ErpResp{
		Segmentation: []*client.Segmentation{
			{
				DateFrom:     time.Now(),
				DateTo:       time.Now(),
				AddressSapId: "12345",
				AdrSegment:   "SegmentA",
				SegmentId:    1,
			},
		},
	}, nil).Once()
	mockClient.On("FetchData", mock.Anything, int64(11), int64(10)).Return(&client.ErpResp{
		Segmentation: []*client.Segmentation{},
	}, nil).Once()
	mockDb.On("NamedExecContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Once()

	ctx := context.Background()
	err := importer.Run(ctx)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func TestImporter_Run_NoMoreData(t *testing.T) {
	mockClient := new(MockClient)
	mockDb := new(MockDb)

	cfg := &config.Config{
		ImportBatchSize:  10,
		ConnInterval:     1000,
		LogCleanupMaxAge: 7,
		LogPath:          "test_log.txt",
	}

	logger := logger.New(cfg)

	importer := New(cfg, logger, mockClient, mockDb)

	mockClient.On("FetchData", mock.Anything, int64(1), int64(10)).Return(&client.ErpResp{
		Segmentation: []*client.Segmentation{},
	}, nil).Once()

	ctx := context.Background()
	err := importer.Run(ctx)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

func TestImporter_Work_InsertError(t *testing.T) {
	mockClient := new(MockClient)
	mockDb := new(MockDb)

	cfg := &config.Config{
		ImportBatchSize:  10,
		ConnInterval:     1000,
		LogCleanupMaxAge: 7,
		LogPath:          "test_log.txt",
	}

	logger := logger.New(cfg)

	importer := New(cfg, logger, mockClient, mockDb)

	mockClient.On("FetchData", mock.Anything, int64(1), int64(10)).Return(&client.ErpResp{
		Segmentation: []*client.Segmentation{
			{
				DateFrom:     time.Now(),
				DateTo:       time.Now(),
				AddressSapId: "12345",
				AdrSegment:   "SegmentA",
				SegmentId:    1,
			},
		},
	}, nil).Once()

	mockDb.On("NamedExecContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("insert error")).Once()

	ctx := context.Background()
	err := importer.work(ctx, 1, 10)

	assert.EqualError(t, err, "insertion failed: insert error")
	mockClient.AssertExpectations(t)
	mockDb.AssertExpectations(t)
}

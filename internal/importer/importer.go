package importer

import (
	"context"
	"errors"
	"fmt"
	"sap_segmentation/internal/client"
	"sap_segmentation/internal/config"
	"sap_segmentation/internal/db"
	"sap_segmentation/internal/formatter"
	"sap_segmentation/internal/logger"
	"sap_segmentation/internal/model"
	"time"
)

var (
	ErrNoMoreData = errors.New("no more data")
)

type Importer struct {
	cfg    *config.Config
	logger logger.Logger
	client client.Client
	db     db.Db
}

func New(cfg *config.Config, logger logger.Logger, client client.Client, db db.Db) *Importer {
	return &Importer{
		cfg:    cfg,
		logger: logger,
		client: client,
		db:     db,
	}
}

func (i *Importer) Run(ctx context.Context) error {
	i.logger.CleanupOldLogs(i.cfg.LogCleanupMaxAge)

	var offset int64 = 1
	limit := i.cfg.ImportBatchSize
	for {
		err := i.work(ctx, offset, limit)
		if err != nil {
			if errors.Is(ErrNoMoreData, err) {
				return nil
			}
			i.logger.Error("importer: %s", err)
		} else {
			offset += limit
		}
		time.Sleep(time.Duration(i.cfg.ConnInterval) * time.Millisecond)
	}
}

func (i *Importer) work(ctx context.Context, offset, limit int64) error {
	data, err := i.client.FetchData(ctx, offset, limit)
	if err != nil {
		return fmt.Errorf("fetch failed: %s", err)
	}
	if len(data.Segmentation) == 0 {
		return ErrNoMoreData
	}
	err = model.InsertSegmentations(
		ctx,
		i.db,
		formatter.ClientToModel(data.Segmentation),
	)
	if err != nil {
		return fmt.Errorf("insertion failed: %s", err)
	}
	return nil
}

package model

import (
	"context"
	"sap_segmentation/internal/db"
)

type Segmentation struct {
	AddressSapId string `db:"address_sap_id"`
	AdrSegment   string `db:"adr_segment"`
	SegmentId    int64  `db:"segment_id"`
}

const (
	SqlInsertSegmentations = `
		INSERT INTO segmentation ("address_sap_id", "adr_segment", "segment_id")
		VALUES (:address_sap_id, :adr_segment, :segment_id)
		ON CONFLICT (address_sap_id) DO UPDATE SET
			adr_segment = EXCLUDED.adr_segment,
			segment_id = EXCLUDED.segment_id
	`
)

func InsertSegmentations(ctx context.Context, db db.Db, segmentations []*Segmentation) error {
	_, err := db.NamedExecContext(ctx, SqlInsertSegmentations, segmentations)
	return err
}

//go:build postgres
// +build postgres

package model

import (
	"context"
	"sap_segmentation/internal/config"
	"sap_segmentation/internal/db"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SegmentationTestSuite struct {
	suite.Suite
	db db.Db
}

func (s *SegmentationTestSuite) SetupSuite() {
	cfg := &config.Config{
		DbHost:     "localhost",
		DbPort:     "5432",
		DbName:     "postgres",
		DbPassword: "admin",
		DbUser:     "admin",
	}
	s.db = db.New(cfg)
}

func (s *SegmentationTestSuite) SetupTest() {
	s.db.GetDb().MustExec("TRUNCATE segmentation")
}

func TestSegmentationTestSuite(t *testing.T) {
	suite.Run(t, new(SegmentationTestSuite))
}

func (s *SegmentationTestSuite) TestInsertSegmentations_Success() {
	segmentations := []*Segmentation{
		{
			AddressSapId: "12345",
			AdrSegment:   "SegmentA",
			SegmentId:    1,
		},
		{
			AddressSapId: "67890",
			AdrSegment:   "SegmentB",
			SegmentId:    2,
		},
	}

	ctx := context.Background()
	err := InsertSegmentations(ctx, s.db, segmentations)

	s.Require().NoError(err)

	db := s.db.GetDb()
	insertedCount := 0
	db.QueryRow("select count(*) from segmentation").Scan(&insertedCount)
	s.Require().Equal(2, insertedCount)
}

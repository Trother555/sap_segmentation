package formatter

import (
	"testing"
	"time"

	"sap_segmentation/internal/client"
	"sap_segmentation/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestClientToModel(t *testing.T) {
	tests := []struct {
		name     string
		input    []*client.Segmentation
		expected []*model.Segmentation
	}{
		{
			name: "Success - Single Segmentation",
			input: []*client.Segmentation{
				{
					DateFrom:     time.Now(),
					DateTo:       time.Now(),
					AddressSapId: "12345",
					AdrSegment:   "SegmentA",
					SegmentId:    1,
				},
			},
			expected: []*model.Segmentation{
				{
					AddressSapId: "12345",
					AdrSegment:   "SegmentA",
					SegmentId:    1,
				},
			},
		},
		{
			name: "Success - Multiple Segmentations",
			input: []*client.Segmentation{
				{
					DateFrom:     time.Now(),
					DateTo:       time.Now(),
					AddressSapId: "12345",
					AdrSegment:   "SegmentA",
					SegmentId:    1,
				},
				{
					DateFrom:     time.Now(),
					DateTo:       time.Now(),
					AddressSapId: "67890",
					AdrSegment:   "SegmentB",
					SegmentId:    2,
				},
			},
			expected: []*model.Segmentation{
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
			},
		},
		{
			name:     "Success - Empty Input",
			input:    []*client.Segmentation{},
			expected: []*model.Segmentation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClientToModel(tt.input)

			assert.Equal(t, tt.expected, result)
		})
	}
}

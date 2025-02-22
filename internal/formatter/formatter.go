package formatter

import (
	"sap_segmentation/internal/client"
	"sap_segmentation/internal/model"
)

func ClientToModel(segmentations []*client.Segmentation) []*model.Segmentation {
	res := make([]*model.Segmentation, len(segmentations))
	for i, seg := range segmentations {
		res[i] = &model.Segmentation{
			AddressSapId: seg.AddressSapId,
			AdrSegment:   seg.AdrSegment,
			SegmentId:    seg.SegmentId,
		}
	}
	return res
}

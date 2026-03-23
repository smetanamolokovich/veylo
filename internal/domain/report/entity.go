package report

import (
	"errors"
	"time"
)

type Report struct {
	id           string
	inspectionID string
	orgID        string
	s3Key        string
	url          string
	generatedAt  time.Time
}

func NewReport(id, inspectionID, orgID, s3Key, url string) (*Report, error) {
	if id == "" || inspectionID == "" || orgID == "" || s3Key == "" {
		return nil, errors.New("report: id, inspectionID, orgID and s3Key are required")
	}
	return &Report{
		id:           id,
		inspectionID: inspectionID,
		orgID:        orgID,
		s3Key:        s3Key,
		url:          url,
		generatedAt:  time.Now().UTC(),
	}, nil
}

func Reconstitute(id, inspectionID, orgID, s3Key, url string, generatedAt time.Time) *Report {
	return &Report{
		id:           id,
		inspectionID: inspectionID,
		orgID:        orgID,
		s3Key:        s3Key,
		url:          url,
		generatedAt:  generatedAt,
	}
}

func (r *Report) ID() string           { return r.id }
func (r *Report) InspectionID() string { return r.inspectionID }
func (r *Report) OrgID() string        { return r.orgID }
func (r *Report) S3Key() string        { return r.s3Key }
func (r *Report) URL() string          { return r.url }
func (r *Report) GeneratedAt() time.Time { return r.generatedAt }

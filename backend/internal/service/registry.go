package service

import "nisst/internal/repository"

type Registry struct {
	Import    *ImportService
	Metadata  *MetadataService
	Analytics *AnalyticsService
	Records   *RecordService
}

func NewRegistry(r *repository.Registry) *Registry {
	return &Registry{
		Import: NewImportService(r), Metadata: NewMetadataService(r),
		Analytics: NewAnalyticsService(r), Records: NewRecordService(r),
	}
}


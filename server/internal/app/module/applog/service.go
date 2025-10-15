package applog

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	"database/sql"
	"errors"
)

type AppLogService struct {
	repo AppLogRepository
}

func NewAppLogService(repo AppLogRepository) *AppLogService {
	return &AppLogService{repo: repo}
}

func (s *AppLogService) GetListAppLogError(dto dto.PaginationQuery) ([]*LogEntry, *response.PaginationMeta, *apperror.AppError) {
	limit := *dto.Limit
	page := *dto.Page

	logs, err := s.repo.ReadPaginatedErrorLogs(int(page), int(limit))
	if err != nil {
		return nil, nil, apperror.Internal("Failed to get applogs", err)
	}

	count, err := s.repo.CountErrorLogEntries()
	if err != nil {
		return nil, nil, apperror.Internal("Failted to get count logs", err)
	}

	pagination := helpers.BuildPaginationMeta(
		limit,
		page,
		int64(count),
		nil,
	)

	return logs, pagination, nil
}

func (s *AppLogService) GetLogByTraceID(traceID string) (*LogEntry, *apperror.AppError) {
	log, err := s.repo.FindLogByTraceID(traceID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("TraceId not found", err)
		}
		return nil, apperror.Internal("Failted to get applog", err)
	}

	return log, nil
}

package applog

type AppLogRepository interface {
	ReadPaginatedErrorLogs(page, limit int) ([]*LogEntry, error)
	CountErrorLogEntries() (int, error)
	FindLogByTraceID(traceID string) (*LogEntry, error)
}

package applog

import (
	"ChainServer/internal/common/config"
	"ChainServer/internal/common/utils"
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type fileAppLogRepository struct {
	root string
}

func NewFileAppLogRepository() AppLogRepository {
	return &fileAppLogRepository{root: config.AppRoot()}
}

func (r *fileAppLogRepository) ReadPaginatedErrorLogs(page, limit int) ([]*LogEntry, error) {
	if page < 1 || limit < 1 {
		return nil, errors.New("invalid page or limit")
	}

	filePath := filepath.Join(r.root, "/logs/error.log")

	if !utils.FileExists(filePath) {
		log.Warn("Applog file is not existing.")
		return make([]*LogEntry, 0), nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var logs []*LogEntry
	rdr := bufio.NewScanner(f)
	skip := (page - 1) * limit
	count := 0
	lineIndex := 0

	for rdr.Scan() {
		if lineIndex < skip {
			lineIndex++
			continue
		}
		var entry LogEntry
		if err := json.Unmarshal(rdr.Bytes(), &entry); err != nil {
			continue
		}
		logs = append(logs, &entry)
		count++
		if count >= limit {
			break
		}
	}

	if err := rdr.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *fileAppLogRepository) CountErrorLogEntries() (int, error) {
	filePath := filepath.Join(r.root, "/logs/error.log")

	if !utils.FileExists(filePath) {
		return 0, nil
	}

	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *fileAppLogRepository) FindLogByTraceID(traceID string) (*LogEntry, error) {
	filepath := filepath.Join(r.root, "/logs/error.log")
	f, err := os.Open(filepath)

	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var log LogEntry

		if err := json.Unmarshal(scanner.Bytes(), &log); err != nil {
			return nil, err
		}

		if log.TraceID == traceID {
			return &log, nil
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return nil, nil
}

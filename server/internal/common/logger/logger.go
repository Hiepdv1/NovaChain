package logger

import (
	"ChainServer/internal/common/config"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type DefaultFieldsHook struct {
	fields log.Fields
}

func (h *DefaultFieldsHook) Levels() []log.Level {
	return log.AllLevels
}

func (h *DefaultFieldsHook) Fire(e *log.Entry) error {
	for k, v := range h.fields {
		if _, exists := e.Data[k]; !exists {
			e.Data[k] = v
		}
	}
	return nil
}

func InitLogger(env string) {
	isProduction := env == "production"
	logLevel := log.DebugLevel

	if isProduction {
		logLevel = log.ErrorLevel
	}

	log.SetLevel(logLevel)

	if isProduction {
		log.SetOutput(colorable.NewColorableStderr())
	} else {
		log.SetOutput(colorable.NewColorableStdout())
	}

	log.SetFormatter(&log.TextFormatter{
		ForceColors:     !isProduction,
		DisableColors:   isProduction,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.AddHook(&DefaultFieldsHook{
		fields: log.Fields{
			"log_scope": "debug",
			"trace_id":  uuid.New().String(),
		},
	})

	dir := filepath.Join(config.AppRoot(), "logs")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Panicf("Failed to create log directory: %v", err)
	}

	if isProduction {
		addRotateHook(dir, log.ErrorLevel, "error")
	} else {
		addRotateHook(dir, log.InfoLevel, "info")
		addRotateHook(dir, log.WarnLevel, "warn")
		addRotateHook(dir, log.ErrorLevel, "error")
	}
}

func addRotateHook(dir string, level log.Level, name string) {
	file := filepath.Join(dir, fmt.Sprintf("%s.log", name))

	hook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   file,
		MaxSize:    2,
		MaxBackups: 7,
		MaxAge:     1,
		Compress:   true,
		Level:      level,
		Formatter: &log.JSONFormatter{
			TimestampFormat: time.RFC3339,
		},
	})

	if err != nil {
		log.Panicf("Failed to initialize rotate hook: %v", err)
	}

	log.AddHook(hook)

}

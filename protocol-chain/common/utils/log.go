package utils

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

var (
	_, file, _, _ = runtime.Caller(0)

	Root = filepath.Join(filepath.Dir(file), "../../")
)

func SetLog(InstanceId string) {
	var logLevel = log.InfoLevel

	dir := path.Join(Root, "/logs/")

	fileName := path.Join("console.log")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	if InstanceId != "" {
		fileName = path.Join(dir, fmt.Sprintf("console_%s.log", InstanceId))
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   fileName,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     30, // days
		Level:      logLevel,
		Compress:   true,
		Formatter: &log.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		log.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.SetLevel(logLevel)
	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	log.AddHook(rotateFileHook)
}

func ClearLogFile(relativePath string) error {
	fullPath := filepath.Join(Root, relativePath)

	absRoot, err := filepath.Abs(Root)
	if err != nil {
		return fmt.Errorf("failed to resolve root dir: %v", err)
	}

	absFull, err := filepath.Abs(fullPath)
	if err != nil {
		return fmt.Errorf("failed to resolve full path: %v", err)
	}

	if len(absFull) < len(absRoot) || absFull[:len(absRoot)] != absRoot {
		return fmt.Errorf("access denied: file outside root dir")
	}

	if _, err := os.Stat(absFull); os.IsNotExist(err) {
		return fmt.Errorf("log file does not exist: %s", absFull)
	}

	file, err := os.OpenFile(absFull, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer file.Close()

	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate log file: %v", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek start: %v", err)
	}

	return nil
}

package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sap_segmentation/internal/config"
	"strings"
	"time"
)

type Logger interface {
	Info(format string, a ...any)
	Debug(format string, a ...any)
	Error(format string, a ...any)
	Warn(format string, a ...any)
	CleanupOldLogs(daysToExpire int64)
}

type LoggerImpl struct {
	l       *log.Logger
	logFile *os.File
	cfg     *config.Config
}

func New(cfg *config.Config) Logger {
	logFile, err := os.Create(cfg.LogPath)
	if err != nil {
		log.Fatalf("failed to open log file: %s", err)
	}
	w := io.MultiWriter(os.Stdout, logFile)
	res := &LoggerImpl{
		l:       log.New(w, "", log.LstdFlags),
		logFile: logFile,
		cfg:     cfg,
	}
	return res
}

func (l *LoggerImpl) CleanupOldLogs(daysToExpire int64) {
	logsDir := filepath.Dir(l.cfg.LogPath)
	files, err := os.ReadDir(logsDir)
	if err != nil {
		l.Warn("failed to cleanup logs: %s", err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		info, err := f.Info()
		if err != nil {
			l.Warn("failed to get file info: %s", err)
			continue
		}

		if strings.HasSuffix(info.Name(), ".log") &&
			info.ModTime().Before(time.Now().Add(-time.Duration(daysToExpire)*24*time.Hour)) {
			err = os.Remove(info.Name())
			l.Warn("failed to remove log file: %s", err)
		}
	}
}

func (l *LoggerImpl) Flush() {
	l.logFile.Sync()
}

func (l *LoggerImpl) Info(format string, a ...any) {
	l.l.Printf("%s: %s", "[INFO]", fmt.Sprintf(format, a...))
}

func (l *LoggerImpl) Debug(format string, a ...any) {
	l.l.Printf("%s: %s", "[DEBUG]", fmt.Sprintf(format, a...))
}

func (l *LoggerImpl) Error(format string, a ...any) {
	l.l.Printf("%s: %s", "[ERROR]", fmt.Sprintf(format, a...))
}

func (l *LoggerImpl) Warn(format string, a ...any) {
	l.l.Printf("%s: %s", "[WARN]", fmt.Sprintf(format, a...))
}

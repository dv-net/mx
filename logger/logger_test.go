package logger_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/dv-net/mx/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogger(t *testing.T) {
	l := logger.New(
		logger.WithAppName("some app name"),
		logger.WithAppVersion("v0.1.0"),
		logger.WithLogLevel(logger.LogLevelDebug),
		logger.WithCaller(true),
		logger.WithStackTrace(true),
		logger.WithZapOption(zap.Hooks(func(entry zapcore.Entry) error {
			fmt.Println("hook")
			return nil
		})),
		logger.WithMemoryBuffer(100),
	)

	l.Info("Hello world")
}

func Test_LoggerWith(t *testing.T) {
	l := logger.New(
		logger.WithAppName("test app name"),
		logger.WithLogLevel(logger.LogLevelDebug),
	)

	l = logger.With(l, "key", "value")

	l = logger.With(l, "key2", "value2")

	l = logger.With(l, "key3", "value3")

	l.Infow("some test value", "numbers", 1234)
}

func TestLogger_WithLogBuffer(t *testing.T) {
	// log without buffer
	l1 := logger.NewExtended(
		logger.WithAppName("no-buffer"),
		logger.WithLogLevel(logger.LogLevelDebug),
	)

	if l1.LastLogs() != nil {
		t.Fatalf("Lastlogs must return Nil if the buffer is not turned on")
	}

	// log with buffer size 3
	l2 := logger.NewExtended(
		logger.WithAppName("with-buffer"),
		logger.WithLogLevel(logger.LogLevelDebug),
		logger.WithMemoryBuffer(3),
	)

	l2.Info("first")
	l2.Debug("second")
	l2.Error("third")
	l2.Warn("fourth")

	logs := l2.LastLogs()
	if len(logs) != 3 {
		t.Fatalf("the buffer should have a maximum of 3 records, got %d", len(logs))
	}

	if logs[0].Level != "warn" {
		t.Errorf("expected level 'warn', got %q", logs[0].Level)
	}
	if logs[0].Message != "fourth" {
		t.Errorf("expected message 'fourth', got %q", logs[0].Message)
	}

	if logs[1].Message != "third" {
		t.Errorf("expected message 'third', got %q", logs[1].Message)
	}
	if logs[2].Message != "second" {
		t.Errorf("expected message 'second', got %q", logs[2].Message)
	}

	if !logs[0].Time.After(time.Now().Add(-time.Second)) {
		t.Errorf("log timestamp is too old: %v", logs[0].Time)
	}
	if !logs[0].Time.Before(time.Now().Add(time.Second)) {
		t.Errorf("log timestamp is too new: %v", logs[0].Time)
	}
}

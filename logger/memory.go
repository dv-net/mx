package logger

import (
	"sync"
	"time"

	"go.uber.org/zap/zapcore"
)

type MemoryLog struct {
	Time    time.Time `json:"time"`
	Level   string    `json:"level"`
	Message string    `json:"message"`
}

type memoryBuffer struct {
	mu   sync.Mutex
	logs []MemoryLog
}
type memoryCore struct {
	inner    zapcore.Core
	buffer   *memoryBuffer
	capacity int
}

func (mem *memoryCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	mem.buffer.mu.Lock()
	defer mem.buffer.mu.Unlock()

	lg := MemoryLog{
		Time:    ent.Time,
		Level:   ent.Level.String(),
		Message: ent.Message,
	}
	mem.buffer.logs = append([]MemoryLog{lg}, mem.buffer.logs...)
	if len(mem.buffer.logs) > mem.capacity {
		mem.buffer.logs = mem.buffer.logs[:mem.capacity]
	}
	return mem.inner.Write(ent, fields)
}

func (mem *memoryCore) With(fields []zapcore.Field) zapcore.Core {
	return &memoryCore{
		inner:    mem.inner.With(fields),
		buffer:   mem.buffer,
		capacity: mem.capacity,
	}
}

func (mem *memoryCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if mem.Enabled(ent.Level) {
		return ce.AddCore(ent, mem)
	}
	return ce
}

func (mem *memoryCore) Sync() error {
	return mem.inner.Sync()
}

func (mem *memoryCore) Enabled(level zapcore.Level) bool {
	return mem.inner.Enabled(level)
}

func (mem *memoryCore) LastLogs() []MemoryLog {
	mem.buffer.mu.Lock()
	defer mem.buffer.mu.Unlock()
	return append([]MemoryLog(nil), mem.buffer.logs...)
}

func newMemoryCore(capacity int, inner zapcore.Core) *memoryCore {
	return &memoryCore{
		inner:    inner,
		capacity: capacity,
		buffer:   &memoryBuffer{logs: make([]MemoryLog, 0, capacity)},
	}
}

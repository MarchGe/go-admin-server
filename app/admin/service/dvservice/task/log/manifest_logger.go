package log

import (
	"bufio"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"log/slog"
	"os"
	"sync"
)

type ManifestLogger struct {
	file   *os.File
	endMtx sync.RWMutex
	end    bool
}

func NewManifestLogger(file *os.File) *ManifestLogger {
	return &ManifestLogger{
		file: file,
	}
}

func (s *ManifestLogger) Original() *os.File {
	return s.file
}

func (s *ManifestLogger) Append(entry *ManifestEntry) {
	s.endMtx.RLock()
	defer s.endMtx.RUnlock()
	if s.end {
		slog.Error("cannot append to ended manifest log")
		return
	}
	if _, err := s.file.WriteString(entry.FormatString() + constant.NewLine); err != nil {
		slog.Error("append manifest log error", slog.Any("err", err))
	}
}

func (s *ManifestLogger) GetReader() *bufio.Reader {
	return bufio.NewReader(s.file)
}

func (s *ManifestLogger) WriteEnd() {
	s.endMtx.Lock()
	defer s.endMtx.Unlock()
	if s.end {
		return
	}
	if _, err := s.file.WriteString(EndFlag); err != nil {
		slog.Error("write log end flag error", slog.Any("err", err))
	} else {
		s.end = true
	}
}

func (s *ManifestLogger) Close() {
	s.file.Close()
}

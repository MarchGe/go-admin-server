package log

import (
	"bufio"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type HostLogger struct {
	file   *os.File
	endMtx sync.RWMutex
	end    bool
}

func NewHostLogger(file *os.File) *HostLogger {
	return &HostLogger{
		file: file,
	}
}

func (s *HostLogger) Original() *os.File {
	return s.file
}

func (s *HostLogger) Append(lineContent string) {
	s.endMtx.RLock()
	defer s.endMtx.RUnlock()
	if s.end {
		slog.Error("cannot append to ended host log")
		return
	}
	nowTime := time.Now().Format(time.DateTime)
	if _, err := s.file.WriteString(nowTime + " " + lineContent + constant.NewLine); err != nil {
		slog.Error("append host log error", slog.Any("err", err))
	}
}

func (s *HostLogger) GetReader() *bufio.Reader {
	return bufio.NewReader(s.file)
}

func (s *HostLogger) WriteEnd() {
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

func (s *HostLogger) GetName() string {
	return filepath.Base(s.file.Name())
}

func (s *HostLogger) Close() {
	s.file.Close()
}

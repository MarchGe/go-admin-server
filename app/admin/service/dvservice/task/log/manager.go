package log

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/service/dvservice"
	"github.com/google/uuid"
	"log/slog"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

type Manager struct {
	dvservice.WorkDirService
	LocateDir string // workdir下的相对路径
}

func NewManager(locateDir string) *Manager {
	return &Manager{
		LocateDir: locateDir,
	}
}

func (s *Manager) CreateManifestLogger(taskId int64) (*ManifestLogger, error) {
	fileName := "manifest-" + time.Now().Format("2006-01-02+15_04_05")
	filePath := fmt.Sprintf("%s/%d/manifest/%s", s.LocateDir, taskId, fileName)
	f, err := s.CreateFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewManifestLogger(f), nil
}

func (s *Manager) CreateHostLogger(taskId int64, host string) (*HostLogger, error) {
	fileName := strings.ReplaceAll(uuid.NewString(), "-", "")
	filePath := fmt.Sprintf("%s/%d/%s/%s", s.LocateDir, taskId, strings.ReplaceAll(host, ".", "_"), fileName)
	f, err := s.CreateFile(filePath)
	if err != nil {
		return nil, err
	}
	return NewHostLogger(f), nil
}

func (s *Manager) OpenLatestManifestLogger(taskId int64) (*ManifestLogger, error) {
	manifestLogDir := s.getManifestLogDir(taskId)
	entries, err := os.ReadDir(manifestLogDir)
	if err != nil {
		return nil, fmt.Errorf("read manifest log dir error, %w", err)
	}
	if len(entries) > 0 {
		s.sortByNameDesc(entries)
		file, err := os.Open(manifestLogDir + "/" + entries[0].Name())
		return NewManifestLogger(file), err
	}
	return nil, errors.New("no file found: " + entries[0].Name())
}

func (s *Manager) OpenHostLogger(taskId int64, host string, hostLogName string) (*HostLogger, error) {
	dir := s.getHostLogDir(taskId, host)
	file, err := os.Open(dir + "/" + hostLogName)
	if err != nil {
		return nil, err
	}
	return NewHostLogger(file), nil
}

// RemoveOldLogs 保留最近的remainHistoryCount次日志，会清除较老的日志
func (s *Manager) RemoveOldLogs(taskId int64, remainHistoryCount int) {
	manifestLogDir := s.getManifestLogDir(taskId)
	entries, err := os.ReadDir(manifestLogDir)
	if err != nil {
		slog.Error("Read manifest log dir error", slog.String("dir", manifestLogDir), slog.Any("err", err))
		return
	}
	s.sortByNameDesc(entries)
	if len(entries) > remainHistoryCount {
		s.removeLogs(taskId, entries[remainHistoryCount:])
	}
}

func (s *Manager) getTaskLogDir(taskId int64) string {
	return s.GetWorkDir() + "/" + s.LocateDir + "/" + strconv.Itoa(int(taskId))
}

func (s *Manager) getManifestLogDir(taskId int64) string {
	return s.getTaskLogDir(taskId) + "/manifest"
}

func (s *Manager) getHostLogDir(taskId int64, host string) string {
	return s.getTaskLogDir(taskId) + "/" + strings.ReplaceAll(host, ".", "_")
}

func (s *Manager) sortByNameDesc(entries []os.DirEntry) {
	slices.SortFunc(entries, func(a, b os.DirEntry) int {
		return -strings.Compare(a.Name(), b.Name())
	})
}

func (s *Manager) removeLogs(taskId int64, manifestEntries []os.DirEntry) {
	manifestLogDir := s.getManifestLogDir(taskId)
	for _, entry := range manifestEntries {
		manifestFile, err := os.Open(manifestLogDir + "/" + entry.Name())
		if err != nil {
			slog.Error("Open manifest log file error", slog.Any("err", err))
			continue
		}
		hostLogRemoveFailed := false
		scanner := bufio.NewScanner(manifestFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || IsEnd(line) {
				continue
			}
			manifestEntry, err := ParseEntry(line)
			if err != nil {
				slog.Error("Parse manifest entry error", slog.Any("err", err))
				continue
			}
			hostLogDir := s.getHostLogDir(taskId, manifestEntry.Host)
			if err = os.Remove(hostLogDir + "/" + manifestEntry.HostLogName); err != nil && !errors.Is(err, os.ErrNotExist) {
				slog.Error("Remove host log file error", slog.Any("err", err))
				hostLogRemoveFailed = true
			}
			entries, e := os.ReadDir(hostLogDir)
			if e == nil && len(entries) == 0 {
				os.Remove(hostLogDir)
			}
		}
		manifestFile.Close()
		if !hostLogRemoveFailed {
			if err = os.Remove(manifestLogDir + "/" + entry.Name()); err != nil && !errors.Is(err, os.ErrNotExist) {
				slog.Error("Remove manifest log file error", slog.Any("err", err))
			}
		}
	}
}

func (s *Manager) RemoveLogs(taskId int64) {
	dir := s.getTaskLogDir(taskId)
	os.RemoveAll(dir)
}

package task

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/config"
	"log/slog"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
)

var _logService = &LogService{}

const LogEndFlag = "-[EOF]-"

var ErrReachedLogEnd = errors.New("reached log end")

type LogService struct {
}

func GetTaskLogService() *LogService {
	return _logService
}

type ManifestEntry struct {
	Index       int    `json:"index"`       // 序号
	Time        string `json:"time"`        // 时间
	Host        string `json:"host"`        // 主机地址（IP）
	HostLogName string `json:"hostLogName"` // 主机日志名称
	Status      string `json:"status"`      // 执行状态，成功、失败等
}

func (s *LogService) createManifestLogFile(taskId int64) (*os.File, error) {
	dir := s.getManifestLogDir(taskId)
	if err := os.MkdirAll(dir, 666); err != nil {
		return nil, err
	}
	fileName := "manifest-" + time.Now().Format("2006-01-02+15_04_05")
	return os.Create(dir + "/" + fileName)
}

func (s *LogService) createHostLogFile(taskId int64, host, fileName string) (*os.File, error) {
	dir := s.getHostLogDir(taskId, host)
	if err := os.MkdirAll(dir, 666); err != nil {
		return nil, err
	}
	return os.Create(dir + "/" + fileName)
}

const manifestEntryFormat = "%d\t%s\t%s\t%s\t%s" // index time host hostLogName status
func (s *LogService) appendManifestLog(manifestFile *os.File, index int, host, hostFileName, status string) {
	nowTime := time.Now().Format(time.DateTime)
	line := fmt.Sprintf(manifestEntryFormat, index, nowTime, host, hostFileName, status)
	if _, err := manifestFile.WriteString(line + constant.NewLine); err != nil {
		slog.Error("append manifest log error", slog.Any("err", err))
	}
}

func (s *LogService) ParseManifestEntry(entryLine string) (*ManifestEntry, error) {
	if s.IsLogEnd(entryLine) {
		return nil, ErrReachedLogEnd
	}
	elements := strings.Split(entryLine, "\t")
	if len(elements) != len(strings.Split(manifestEntryFormat, "\t")) {
		return nil, errors.New("manifest entry error: '" + entryLine + "'")
	}
	index, err := strconv.Atoi(elements[0])
	if err != nil {
		return nil, errors.New("index of entry line error: " + elements[0])
	}
	entry := &ManifestEntry{
		Index:       index,
		Time:        elements[1],
		Host:        elements[2],
		HostLogName: elements[3],
		Status:      elements[4],
	}
	return entry, nil
}

func (s *LogService) appendHostLog(hostFile *os.File, lineContent string) {
	nowTime := time.Now().Format(time.DateTime)
	if _, err := hostFile.WriteString(nowTime + " " + lineContent + constant.NewLine); err != nil {
		slog.Error("append host log error", slog.Any("err", err))
	}
}

// writeEnd 每个日志文件最后都要写上一个结束标志
func (s *LogService) writeEnd(f *os.File) {
	if _, err := f.WriteString(LogEndFlag); err != nil {
		slog.Error("write log end flag error", slog.Any("err", err))
	}
}

func (s *LogService) getTaskLogDir(taskId int64) string {
	return path.Clean(config.GetConfig().WorkDir) + "/data/log/task/" + strconv.Itoa(int(taskId))
}

func (s *LogService) getManifestLogDir(taskId int64) string {
	return s.getTaskLogDir(taskId) + "/manifest"
}

func (s *LogService) getHostLogDir(taskId int64, host string) string {
	return s.getTaskLogDir(taskId) + "/" + strings.ReplaceAll(host, ".", "_")
}

// remainDeployHistory 保留最近的remainHistoryCount次部署日志，会清除较老的日志
func (s *LogService) removeOldLogs(taskId int64, remainHistoryCount int) {
	manifestLogDir := s.getManifestLogDir(taskId)
	entries, err := os.ReadDir(manifestLogDir)
	if err != nil {
		slog.Error("read dir error", slog.String("dir", manifestLogDir), slog.Any("err", err))
		return
	}
	s.sortByNameDesc(entries)
	if len(entries) > remainHistoryCount {
		s.removeLogs(taskId, entries[remainHistoryCount:])
	}
}

func (s *LogService) sortByNameDesc(entries []os.DirEntry) {
	slices.SortFunc(entries, func(a, b os.DirEntry) int {
		return -strings.Compare(a.Name(), b.Name())
	})
}

func (s *LogService) removeLogs(taskId int64, manifestEntries []os.DirEntry) {
	manifestLogDir := s.getManifestLogDir(taskId)
	for _, entry := range manifestEntries {
		manifestFile, err := os.Open(manifestLogDir + "/" + entry.Name())
		if err != nil {
			slog.Error("open file error", slog.Any("err", err))
			continue
		}
		scanner := bufio.NewScanner(manifestFile)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || s.IsLogEnd(line) {
				continue
			}
			manifestEntry, err := s.ParseManifestEntry(line)
			if err != nil {
				slog.Error("parse manifest entry error", slog.Any("err", err))
				continue
			}
			hostLogDir := s.getHostLogDir(taskId, manifestEntry.Host)
			if err = os.Remove(hostLogDir + "/" + manifestEntry.HostLogName); err != nil && !errors.Is(err, os.ErrNotExist) {
				slog.Error("remove file error", slog.Any("err", err))
			}
			entries, e := os.ReadDir(hostLogDir)
			if e == nil && len(entries) == 0 {
				_ = os.Remove(hostLogDir)
			}
		}
		_ = manifestFile.Close()
		if err = os.Remove(manifestLogDir + "/" + entry.Name()); err != nil && !errors.Is(err, os.ErrNotExist) {
			slog.Error("remove file error", slog.Any("err", err))
		}
	}
}

func (s *LogService) OpenLatestManifestLogFile(taskId int64) (*os.File, error) {
	manifestLogDir := s.getManifestLogDir(taskId)
	entries, err := os.ReadDir(manifestLogDir)
	if err != nil {
		return nil, fmt.Errorf("read dir error, %w", err)
	}
	if len(entries) > 0 {
		s.sortByNameDesc(entries)
		file, err := os.Open(manifestLogDir + "/" + entries[0].Name())
		return file, err
	}
	return nil, errors.New("no file found: " + entries[0].Name())
}

func (s *LogService) IsLogEnd(line string) bool {
	return line == LogEndFlag
}

func (s *LogService) OpenHostLogFile(taskId int64, host string, hostLogName string) (*os.File, error) {
	dir := s.getHostLogDir(taskId, host)
	return os.Open(dir + "/" + hostLogName)
}

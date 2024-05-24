package dvservice

import (
	"errors"
	dvRes "github.com/MarchGe/go-admin-server/app/admin/service/dvservice/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"os"
	"strings"
)

var _explorerService = &ExplorerService{}

type ExplorerService struct {
}

func GetExplorerService() *ExplorerService {
	return _explorerService
}

func (s *ExplorerService) ListEntries(parentDir, keyword string) ([]*dvRes.ExplorerEntry, error) {
	info, err := os.Stat(parentDir)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, E.Message("父目录参数有误")
	}
	dirEntries, err := os.ReadDir(parentDir)
	if err != nil {
		if errors.Is(err, os.ErrPermission) {
			return nil, E.Message("没有权限访问该目录")
		}
		return nil, err
	}
	var length = 0
	if keyword == "" {
		length = len(dirEntries)
	}
	for _, item := range dirEntries {
		if keyword != "" && strings.Contains(item.Name(), keyword) {
			length++
		}
	}
	entries := make([]*dvRes.ExplorerEntry, length)
	for i, item := range dirEntries {
		if keyword == "" || keyword != "" && strings.Contains(item.Name(), keyword) {
			entry := &dvRes.ExplorerEntry{
				Name: item.Name(),
				Type: s.parseType(item.Type()),
			}
			entries[i] = entry
		}
	}
	return entries, nil
}

func (s *ExplorerService) parseType(mode os.FileMode) dvRes.EntryType {
	fileType := mode.Type()
	switch {
	case fileType.IsDir():
		return dvRes.EntryTypeDir
	case fileType&os.ModeSymlink == os.ModeSymlink:
		return dvRes.EntryTypeLink
	case fileType&os.ModeSocket == os.ModeSocket:
		return dvRes.EntryTypeSocket
	case fileType&os.ModeNamedPipe == os.ModeNamedPipe:
		return dvRes.EntryTypeNamedPipe
	case fileType&os.ModeDevice == os.ModeDevice:
		if fileType&os.ModeCharDevice == os.ModeCharDevice {
			return dvRes.EntryTypeCharDevice
		} else {
			return dvRes.EntryTypeBlockDevice
		}
	}
	return dvRes.EntryTypeDefault
}

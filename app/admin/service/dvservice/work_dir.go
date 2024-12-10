package dvservice

import (
	"github.com/MarchGe/go-admin-server/config"
	"os"
	"path"
)

var _workDirService = &WorkDirService{}

type WorkDirService struct {
	__workDir string
}

func GetWorkdirService() *WorkDirService {
	return _workDirService
}

func (s *WorkDirService) GetWorkDir() string {
	if s.__workDir == "" {
		s.__workDir = path.Clean(config.GetConfig().WorkDir)
	}
	return s.__workDir
}

func (s *WorkDirService) CreateDirectory(relativeDirPath string) error {
	dir := path.Clean(relativeDirPath)
	if path.IsAbs(dir) {
		return os.MkdirAll(s.GetWorkDir()+dir, 0755)
	}
	return os.MkdirAll(s.GetWorkDir()+"/"+dir, 0755)
}

func (s *WorkDirService) RemoveDirectory(relativeDirPath string) error {
	dir := path.Clean(relativeDirPath)
	if path.IsAbs(dir) {
		return os.RemoveAll(s.GetWorkDir() + dir)
	}
	return os.RemoveAll(s.GetWorkDir() + "/" + dir)
}

func (s *WorkDirService) CreateFile(relativeFilePath string) (*os.File, error) {
	file := path.Clean(relativeFilePath)
	absolutePath := s.GetWorkDir() + "/" + file
	if path.IsAbs(file) {
		absolutePath = s.GetWorkDir() + file
	}
	dir := path.Dir(absolutePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return os.Create(absolutePath)
}

func (s *WorkDirService) RemoveFile(relativeFilePath string) error {
	file := path.Clean(relativeFilePath)
	if path.IsAbs(file) {
		return os.Remove(s.GetWorkDir() + file)
	}
	return os.Remove(s.GetWorkDir() + "/" + file)
}

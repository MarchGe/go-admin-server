package log

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const manifestEntryFormat = "%d\t%s\t%s\t%s\t%s" // index time host hostLogName status
const columns = 5

type ManifestEntry struct {
	Index       int    `json:"index"`       // 序号
	Time        string `json:"time"`        // 时间
	Host        string `json:"host"`        // 主机地址（IP）
	HostLogName string `json:"hostLogName"` // 主机日志名称
	Status      string `json:"status"`      // 执行状态，成功、失败等
}

func NewEntry(index int, host, hostLogFileName, status string) *ManifestEntry {
	return &ManifestEntry{
		Index:       index,
		Time:        time.Now().Format(time.DateTime),
		Host:        host,
		HostLogName: hostLogFileName,
		Status:      strings.ToUpper(status),
	}
}

func ParseEntry(entryLine string) (*ManifestEntry, error) {
	elements := strings.Split(entryLine, "\t")
	if len(elements) != columns {
		return nil, errors.New("parse manifest entry error: columns count not match, entry line: '" + entryLine + "'")
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

func (s *ManifestEntry) FormatString() string {
	return fmt.Sprintf(manifestEntryFormat, s.Index, s.Time, s.Host, s.HostLogName, s.Status)
}

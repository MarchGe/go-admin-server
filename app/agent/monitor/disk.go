package monitor

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/model"
	"github.com/shirou/gopsutil/v3/disk"
)

func GetDiskStats() (*model.DiskStat, error) {
	stat, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("get system disk stats error, %w", err)
	}
	stats := &model.DiskStat{
		Total: float32(int64(float32(stat.Total)*100/1024/1024/1024)) / 100,
		Used:  float32(int64(float32(stat.Used)*100/1024/1024/1024)) / 100,
	}
	stats.UsedPercent = float32(int64(stat.UsedPercent*100)) / 100
	return stats, nil
}

package monitor

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/model"
	"github.com/shirou/gopsutil/v3/cpu"
)

func GetCpuStat() (*model.CpuStat, error) {
	physicalCounts, e1 := cpu.Counts(false)
	logicalCounts, e2 := cpu.Counts(true)
	percent, e3 := cpu.Percent(0, false)
	if err := errors.Join(e1, e2, e3); err != nil {
		return nil, fmt.Errorf("get system cpu stats error, %w", err)
	}
	stat := &model.CpuStat{
		PhysicalCores: int32(physicalCounts),
		LogicalCores:  int32(logicalCounts),
		UsedPercent:   float32(int64(percent[0]*100)) / 100,
	}
	return stat, nil
}

func GetCpuInfo() ([]*model.CpuInfo, error) {
	infos, err := cpu.Info()
	if err != nil {
		return nil, fmt.Errorf("get cpu info error, %w", err)
	}
	cpuInfos := make([]*model.CpuInfo, len(infos))
	for i := range infos {
		cpuInfos[i] = &model.CpuInfo{}
		cpuInfos[i].Num = infos[i].CPU
		cpuInfos[i].VendorId = infos[i].VendorID
		cpuInfos[i].Family = infos[i].Family
		cpuInfos[i].PhysicalId = infos[i].PhysicalID
		cpuInfos[i].Cores = infos[i].Cores
		cpuInfos[i].ModelName = infos[i].ModelName
	}
	return cpuInfos, nil
}

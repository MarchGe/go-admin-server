package monitor

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/model"
	"github.com/shirou/gopsutil/v3/mem"
	"strconv"
)

func GetMemoryStat() (virtualMemory *model.MemoryStat, swapMemory *model.MemoryStat, err error) {
	vm, e1 := mem.VirtualMemory()
	sm, e2 := mem.SwapMemory()
	if e := errors.Join(e1, e2); e != nil {
		err = fmt.Errorf("get system memory stats error, %w", e)
		return
	}
	virtualStat := &model.MemoryStat{
		Total: float32(int64(float32(vm.Total)*100/1024/1024/1024)) / 100,
		Used:  float32(int64(float32(vm.Used)*100/1024/1024/1024)) / 100,
	}
	var formatResult string
	if virtualStat.Total != 0 {
		formatResult = fmt.Sprintf("%.2f", virtualStat.Used*100/virtualStat.Total)
	} else {
		formatResult = "0"
	}
	percent, _ := strconv.ParseFloat(formatResult, 32)
	virtualStat.UsedPercent = float32(percent)

	swapStat := &model.MemoryStat{
		Total: float32(int64(float32(sm.Total)*100/1024/1024/1024)) / 100,
		Used:  float32(int64(float32(sm.Used)*100/1024/1024/1024)) / 100,
	}
	if swapStat.Total != 0 {
		formatResult = fmt.Sprintf("%.2f", swapStat.Used*100/swapStat.Total)
	} else {
		formatResult = "0"
	}
	percent, _ = strconv.ParseFloat(formatResult, 32)
	swapStat.UsedPercent = float32(percent)
	return virtualStat, swapStat, nil
}

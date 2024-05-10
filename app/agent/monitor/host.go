package monitor

import (
	"fmt"
	"github.com/MarchGe/go-admin-server/agent/grpc/pb/model"
	"github.com/shirou/gopsutil/v3/host"
	"strconv"
	"strings"
)

func GetHostInfo() (*model.HostInfo, error) {
	info, err := host.Info()
	if err != nil {
		return nil, fmt.Errorf("get host info error, %w", err)
	}
	ip, e := GetIP()
	if e != nil {
		return nil, fmt.Errorf("get machine ip error, %w", e)
	}
	hostInfo := &model.HostInfo{
		Ip:              ip,
		HostName:        info.Hostname,
		UpTime:          transferTime(int64(info.Uptime)),
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
	}
	infos, e := GetCpuInfo()
	if e != nil {
		return nil, fmt.Errorf("GetCpuInfo error, %w", e)
	}
	hostInfo.CpuInfos = infos
	return hostInfo, nil
}

func transferTime(timeInSeconds int64) string {
	var hour, day, month, year string
	var time int64
	if time = timeInSeconds / 3600; time > 0 {
		hour = strconv.Itoa(int(time%24)) + "小时"
	}
	if time = time / 24; time > 0 {
		day = strconv.Itoa(int(time%30)) + "天"
	}
	if time = time / 30; time > 0 {
		month = strconv.Itoa(int(time%12)) + "月"
	}
	if time = time / 12; time > 0 {
		year = strconv.Itoa(int(time)) + "年"
	}
	return strings.Join([]string{year, month, day, hour}, "")
}

package monitor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/admin/grpc/pb/model"
	"github.com/MarchGe/go-admin-server/app/admin/service/dto/res"
	"github.com/MarchGe/go-admin-server/app/common/E"
	"github.com/MarchGe/go-admin-server/app/common/R"
	"github.com/MarchGe/go-admin-server/app/common/constant"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	ginUtils "github.com/MarchGe/go-admin-server/app/common/utils/gin_utils"
	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"
)

var _sysStatsApi = &SysStatsApi{}
var nodeList = make([]string, 0)

type SysStatsApi struct {
}

func GetSysStatsApi() *SysStatsApi {
	return _sysStatsApi
}

// GetList godoc
//
//	@Summary	查询监控服务器列表
//	@Tags		性能监控
//	@Produce	application/json
//	@Param		keyword		query		string	false	"按照IP模糊搜索"
//	@Param		page		query		int64	false	"页码"
//	@Param		pageSize	query		int64	false	"每页查询条数"
//	@Success	200			{object}	R.Result{value=res.PageableData[SysStats]}
//	@Router		/monitor/list [get]
func (s *SysStatsApi) GetList(c *gin.Context) {
	keyword := ginUtils.GetStringQuery(c, "keyword", "")
	page, err1 := ginUtils.GetIntQuery(c, "page", 1)
	pageSize, err2 := ginUtils.GetIntQuery(c, "pageSize", 10)
	err := errors.Join(err1, err2)
	if err != nil {
		R.Fail(c, err.Error(), http.StatusBadRequest)
		return
	}
	if pageSize*(page-1) >= len(nodeList) {
		r := &res.PageableData[*SysStats]{
			List:  make([]*SysStats, 0),
			Total: 0,
		}
		R.Success(c, r)
		return
	}
	selectedIps := make([]string, 0, pageSize)
	skip := 0
	for _, item := range nodeList {
		if strings.HasPrefix(item, keyword) {
			if skip >= pageSize*(page-1) && skip < pageSize*page {
				selectedIps = append(selectedIps, item)
			}
			skip++
		}
	}
	sysStats := make([]*SysStats, len(selectedIps))
	if len(selectedIps) > 0 {
		cache := utils.GetCache()
		for i := range selectedIps {
			bytes, e := cache.Get(constant.CacheKeyPerformanceStats(selectedIps[i]))
			if e != nil {
				if errors.Is(err, bigcache.ErrEntryNotFound) {
					R.Fail(c, "目标主机'"+selectedIps[i]+"'性能数据不存在", http.StatusBadRequest)
					return
				}
				E.PanicErr(e)
			} else {
				var stats = &model.SysStats{}
				if err = json.Unmarshal(bytes, stats); err != nil {
					E.PanicErr(err)
				}
				sysStats[i] = transferRpcMessage(stats)
			}
		}
	}
	result := &res.PageableData[*SysStats]{
		List:  sysStats,
		Total: int64(skip),
	}
	R.Success(c, result)
}

// GetPerformanceStats godoc
//
//	@Summary	获取最新性能监控数据
//	@Tags		性能监控
//	@Produce	application/json
//	@Param		ip	query		string	true	"目标主机ip地址"
//	@Success	200	{object}	R.Result{value=SysStats}
//	@Router		/monitor/performance-stats [get]
func (s *SysStatsApi) GetPerformanceStats(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		R.Fail(c, "目标主机ip地址不能为空", http.StatusBadRequest)
		return
	}
	cache := utils.GetCache()
	bytes, err := cache.Get(constant.CacheKeyPerformanceStats(ip))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			R.Fail(c, "目标主机性能数据不存在", http.StatusBadRequest)
			return
		}
		E.PanicErr(err)
	}
	var stats = &model.SysStats{}
	if err = json.Unmarshal(bytes, stats); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, transferRpcMessage(stats))
}

// GetHostInfo godoc
//
//	@Summary	获取主机信息
//	@Tags		性能监控
//	@Produce	application/json
//	@Param		ip	query		string	true	"目标主机ip地址"
//	@Success	200	{object}	R.Result{value=HostInfo}
//	@Router		/monitor/host-info [get]
func (s *SysStatsApi) GetHostInfo(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		R.Fail(c, "目标主机ip地址不能为空", http.StatusBadRequest)
		return
	}
	cache := utils.GetCache()
	bytes, err := cache.Get(constant.CacheKeyHostInfo(ip))
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			R.Fail(c, "目标主机信息不存在", http.StatusBadRequest)
			return
		}
		E.PanicErr(err)
	}
	var hostInfo = &model.HostInfo{}
	if err = json.Unmarshal(bytes, hostInfo); err != nil {
		E.PanicErr(err)
	}
	R.Success(c, transferRpcHostInfo(hostInfo))
}

// DeleteHost godoc
//
//	@Summary	删除一条主机监控记录
//	@Tags		性能监控
//	@Produce	application/json
//	@Param		ip	query		string	true	"目标主机ip地址"
//	@Success	200	{object}	R.Result
//	@Router		/monitor/host [delete]
func (s *SysStatsApi) DeleteHost(c *gin.Context) {
	ip := c.Query("ip")
	if ip == "" {
		R.Fail(c, "目标主机ip地址不能为空", http.StatusBadRequest)
		return
	}
	mtx.Lock()
	defer mtx.Unlock()
	for i := 0; i < len(nodeList); i++ {
		if nodeList[i] == ip {
			if i != len(nodeList)-1 {
				copy(nodeList[i:], nodeList[i+1:])
			}
			nodeList = nodeList[:len(nodeList)-1]
			break
		}
	}
	cache := utils.GetCache()
	err := cache.Delete(constant.CacheKeyHostInfo(ip))
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		E.PanicErr(err)
	}
	err = cache.Delete(constant.CacheKeyPerformanceStats(ip))
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		E.PanicErr(err)
	}
	R.Success(c, nil)
}

func ProcessHostInformation(hostInfo *model.HostInfo) error {
	cache := utils.GetCache()
	bytes, err := json.Marshal(hostInfo)
	if err != nil {
		return fmt.Errorf("marshal rpc hostInfo message to json error, %w", err)
	}
	if err = cache.Set(constant.CacheKeyHostInfo(hostInfo.Ip), bytes); err != nil {
		return fmt.Errorf("save rpc hostInfo message to cache error, %w", err)
	}
	return nil
}

var mtx sync.Mutex

func ProcessPerformanceStats(stats *model.SysStats) error {
	bytes, err := json.Marshal(stats)
	if err != nil {
		return fmt.Errorf("transfer rpc performanceStats message to json error, %w", err)
	}
	if !slices.Contains(nodeList, stats.Ip) {
		mtx.Lock()
		nodeList = append(nodeList, stats.Ip)
		mtx.Unlock()
	}
	cache := utils.GetCache()
	if err = cache.Set(constant.CacheKeyPerformanceStats(stats.Ip), bytes); err != nil {
		return fmt.Errorf("save rpc message to local cache error, %w", err)
	}
	return nil
}

func transferRpcMessage(stats *model.SysStats) *SysStats { // rpc message传输时会把零值丢弃，所以需要转一下
	t := time.Unix(stats.Timestamp.Seconds, int64(stats.Timestamp.Nanos))
	return &SysStats{
		Ip: stats.Ip,
		Cpu: &CpuStat{
			PhysicalCores: stats.Cpu.PhysicalCores,
			LogicalCores:  stats.Cpu.LogicalCores,
			UsedPercent:   stats.Cpu.UsedPercent,
		},
		VirtualMemory: &MemoryStat{
			Total:       stats.VirtualMemory.Total,
			Used:        stats.VirtualMemory.Used,
			UsedPercent: stats.VirtualMemory.UsedPercent,
		},
		SwapMemory: &MemoryStat{
			Total:       stats.SwapMemory.Total,
			Used:        stats.SwapMemory.Used,
			UsedPercent: stats.SwapMemory.UsedPercent,
		},
		Disk: &DiskStat{
			Total:       stats.Disk.Total,
			Used:        stats.Disk.Used,
			UsedPercent: stats.Disk.UsedPercent,
		},
		Time: &t,
	}
}

func transferRpcHostInfo(info *model.HostInfo) *HostInfo {
	t := time.Unix(info.Timestamp.Seconds, int64(info.Timestamp.Nanos))
	hostInfo := &HostInfo{
		Ip:              info.Ip,
		HostName:        info.HostName,
		UpTime:          info.UpTime,
		Platform:        info.Platform,
		PlatformVersion: info.PlatformVersion,
		KernelVersion:   info.KernelVersion,
		KernelArch:      info.KernelArch,
		Time:            &t,
	}
	cpuInfos := make([]*CpuInfo, len(info.CpuInfos))
	infos := info.CpuInfos
	for i := range infos {
		cpuInfos[i] = &CpuInfo{
			Num:        infos[i].Num,
			VendorId:   infos[i].VendorId,
			Family:     infos[i].Family,
			PhysicalId: infos[i].PhysicalId,
			Cores:      infos[i].Cores,
			ModelName:  infos[i].ModelName,
			Mhz:        infos[i].Mhz,
		}
	}
	hostInfo.CpuInfos = cpuInfos
	return hostInfo
}

type CpuStat struct {
	PhysicalCores int32   `json:"physicalCores"`
	LogicalCores  int32   `json:"logicalCores"`
	UsedPercent   float32 `json:"usedPercent"`
}

type MemoryStat struct {
	Total       float32 `json:"total"`
	Used        float32 `json:"used"`
	UsedPercent float32 `json:"usedPercent"`
}

type DiskStat struct {
	Total       float32 `json:"total"`
	Used        float32 `json:"used"`
	UsedPercent float32 `json:"usedPercent"`
}

type SysStats struct {
	Ip            string      `json:"ip"`
	Cpu           *CpuStat    `json:"cpu"`
	VirtualMemory *MemoryStat `json:"virtualMemory"`
	SwapMemory    *MemoryStat `json:"swapMemory"`
	Disk          *DiskStat   `json:"disk"`
	Time          *time.Time  `json:"time"`
}

type HostInfo struct {
	Ip              string     `json:"ip"`
	HostName        string     `json:"hostName"`
	UpTime          string     `json:"upTime"`
	Platform        string     `json:"platform"`
	PlatformVersion string     `json:"platformVersion"`
	KernelVersion   string     `json:"kernelVersion"`
	KernelArch      string     `json:"kernelArch"`
	CpuInfos        []*CpuInfo `json:"cpuInfos"`
	Time            *time.Time `json:"time"`
}

type CpuInfo struct {
	Num        int32   `json:"num"`
	VendorId   string  `json:"vendorId"`
	Family     string  `json:"family"`
	PhysicalId string  `json:"physicalId"`
	Cores      int32   `json:"cores"`
	ModelName  string  `json:"modelName"`
	Mhz        float32 `json:"mhz"`
}

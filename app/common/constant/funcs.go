package constant

func CacheKeyHostInfo(ip string) string {
	return "h.inf:" + ip
}

func CacheKeyPerformanceStats(ip string) string {
	return "p.s:" + ip
}

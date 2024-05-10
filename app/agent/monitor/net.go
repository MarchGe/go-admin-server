package monitor

import (
	"fmt"
	"log/slog"
	"net"
	"time"
)

var _ip string

func init() {
	tick := time.Tick(30 * time.Second)
	go func() {
		slog.Info("starting time ticker for getting machine ip, interval: 30 seconds")
		for {
			<-tick
			_ip, _ = getIp()
		}
	}()
}

func GetIP() (string, error) {
	if _ip != "" {
		return _ip, nil
	}
	return getIp()
}

func getIp() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("get local ip error, %w", err)
	}
	for i := range interfaces {
		addrs, e := interfaces[i].Addrs()
		if e != nil {
			continue
		}
		for i := range addrs {
			ipNet, ok := addrs[i].(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
				continue
			}
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("cannot find machine ip address")
}

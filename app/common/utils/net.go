package utils

import (
	"fmt"
	"net"
)

func GetIP() (string, error) {
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

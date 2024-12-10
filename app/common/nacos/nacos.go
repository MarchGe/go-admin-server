package nacos

import (
	"errors"
	"fmt"
	"github.com/MarchGe/go-admin-server/app/common/utils"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"log"
	"net"
	"strconv"
)

type Nacoser struct {
	cfg          *Config
	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient
	serviceHost  string
	servicePort  uint64
}

func CreateNacoser(cfg *Config) *Nacoser {
	clientConfig := getClientConfig(cfg)
	serverConfigs := getServerConfigs(cfg.Servers)
	nacosClientParam := vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	}
	configClient, err := clients.NewConfigClient(nacosClientParam)
	if err != nil {
		log.Panicf("create nacos config client error: %v", err)
	}

	namingClient, err := clients.NewNamingClient(nacosClientParam)
	if err != nil {
		log.Panicf("create nacos naming client error: %v", err)
	}
	return &Nacoser{
		cfg:          cfg,
		configClient: configClient,
		namingClient: namingClient,
	}
}

func (n *Nacoser) GetConfig() (config string, formatType string) {
	configParam := vo.ConfigParam{
		DataId: n.cfg.DataId,
		Group:  n.cfg.Group,
		Type:   n.cfg.Type,
		Tag:    n.cfg.Tag,
	}
	stringConfig, err := n.configClient.GetConfig(configParam)
	if err != nil {
		log.Panicf("nacos get config error: %v", err)
	}
	return stringConfig, n.cfg.Type
}

func (n *Nacoser) RegisterService(listenAddr string) error { // 备注：非nacos超级用户向public默认命名空间中注册临时服务，会报授权失败，可能是nacos-sdk-go v2.2.5 或 nacos-server v2.1.2中有bug
	n.parseHostAndPort(listenAddr)
	serviceInfo := n.cfg.ServiceInfo
	ok, err := n.namingClient.RegisterInstance(
		vo.RegisterInstanceParam{
			ServiceName: serviceInfo.ServiceName,
			Ip:          n.serviceHost,
			Port:        n.servicePort,
			Weight:      serviceInfo.Weight,
			Enable:      serviceInfo.Enable,
			Healthy:     serviceInfo.Healthy,
			Ephemeral:   serviceInfo.Ephemeral,
			ClusterName: n.cfg.ClusterName,
			GroupName:   n.cfg.Group,
		},
	)
	if err != nil {
		return fmt.Errorf("nacos naming client register instance error, %w", err)
	}
	if !ok {
		return errors.New("register service to nacos failed, unknown reason")
	}
	return nil
}

func (n *Nacoser) UnregisterService() {
	serviceInfo := n.cfg.ServiceInfo
	if !serviceInfo.Ephemeral {
		ok, err := n.namingClient.DeregisterInstance(
			vo.DeregisterInstanceParam{
				ServiceName: serviceInfo.ServiceName,
				Ip:          n.serviceHost,
				Port:        n.servicePort,
				Cluster:     n.cfg.ClusterName,
				GroupName:   n.cfg.Group,
				Ephemeral:   serviceInfo.Ephemeral,
			},
		)
		if err != nil {
			log.Printf("unregister service '%s' from nacos failed: %v", n.cfg.ServiceInfo.ServiceName, err)
			return
		}
		if !ok {
			log.Printf("unregister service '%s' from nacos failed, unknown reason", n.cfg.ServiceInfo.ServiceName)
			return
		}
		log.Printf("unregister service '%s' success", n.cfg.ServiceInfo.ServiceName)
	}
}

func (n *Nacoser) GetService(param vo.GetServiceParam) *model.Service {
	service, err := n.namingClient.GetService(param)
	if err != nil {
		log.Printf("get service error: %v", err)
		return nil
	}
	return &service
}

func (n *Nacoser) GetServiceInstance(param vo.SelectOneHealthInstanceParam) *model.Instance {
	instance, err := n.namingClient.SelectOneHealthyInstance(param)
	if err != nil {
		log.Printf("get service instance error: %v", err)
		return nil
	}
	return instance
}

func getServerConfigs(sConfigs []ServerConfig) []constant.ServerConfig {
	configs := make([]constant.ServerConfig, len(sConfigs))
	for i, sConfig := range sConfigs {
		configs[i] = constant.ServerConfig{
			IpAddr:   sConfig.IpAddr,
			Port:     sConfig.Port,
			GrpcPort: sConfig.GrpcPort,
		}
	}
	return configs
}

func getClientConfig(cfg *Config) *constant.ClientConfig {
	return &constant.ClientConfig{
		TimeoutMs:           cfg.TimeoutMs,
		NamespaceId:         cfg.NamespaceId,
		Username:            cfg.Username,
		Password:            cfg.Password,
		NotLoadCacheAtStart: cfg.NotLoadCacheAtStart,
		LogLevel:            cfg.LogLevel,
	}
}

func (n *Nacoser) parseHostAndPort(listenAddr string) {
	if n.serviceHost != "" || n.servicePort != 0 {
		return
	}
	host, sPort, err := net.SplitHostPort(listenAddr)
	if err != nil {
		log.Panicf("parse port from 'Listen' address error: %v", err)
	}
	port, err := strconv.Atoi(sPort)
	if err != nil {
		log.Panicf("listen port error: %v", err)
	}
	ip := net.ParseIP(host) // net.ParseIP()只能解析ip格式的参数，如果传的是域名，返回也是nil
	if host == "" || (ip != nil && (ip.IsLoopback() || ip.To4() == nil)) {
		var e error
		host, e = utils.GetIP()
		if e != nil {
			log.Panicf("get machine ip error, %v", e)
		}
		if host == "" {
			log.Panicf("cannot get available ip for registering service")
		}
	}
	n.serviceHost = host
	n.servicePort = uint64(port)
}

func (n *Nacoser) Close() {
	n.configClient.CloseClient()
	n.namingClient.CloseClient()
}

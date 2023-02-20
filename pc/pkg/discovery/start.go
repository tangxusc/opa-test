package discovery

import (
	"context"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/sirupsen/logrus"
	"net"
	"plugin_center/pkg/config"
	"strconv"
	"time"
)

var client naming_client.INamingClient

func Connect(ctx context.Context) error {
	configs := make([]constant.ServerConfig, len(servers))
	path := constant.WithContextPath("/nacos")
	for i, server := range servers {
		configs[i] = *constant.NewServerConfig(server, 8848, path)
	}
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId(Namespace),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel("debug"),
	)
	var err error
	client, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: configs,
		},
	)
	if err != nil {
		return err
	}

	go func() {
		select {
		case <-ctx.Done():
			client.CloseClient()
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				err = register(client)
				if err == nil {
					ticker.Stop()
					return
				} else {
					logrus.Errorf("[discovery]register service failed,error:%v,3s after retry...", err)
				}
			}
		}
	}()
	return nil
}

func register(client naming_client.INamingClient) error {
	atoi, err := strconv.Atoi(config.WebPort)
	if err != nil {
		logrus.Errorf("[discovery]register service failed,convert webport to int error:%v", err)
		return err
	}
	ip, err := getClientIp()
	if err != nil {
		logrus.Errorf("[discovery]get client ip error:%v", err)
		return err
	}
	param := vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(atoi),
		Enable:      true,
		Healthy:     true,
		ServiceName: serviceName,
		Ephemeral:   true,
	}
	success, err := client.RegisterInstance(param)
	if !success || err != nil {
		logrus.Errorf("[discovery]register service failed,error:%v", err)
		return err
	} else {
		logrus.Debugf("[discovery]register service success")
	}

	return nil
}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), err
			}
		}
	}
	return "", err
}

func GetDiscoveryClient() naming_client.INamingClient {
	return client
}

package configs

import (
	"os"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
)

const (
	NacosDefaultGroup = "DEFAULT_GROUP"
)

// 以 nacos 作为配置中心数据源
type SourceNacos struct {
	client config_client.IConfigClient
}

var _ ISource = (*SourceNacos)(nil)

const (
	EnvNacosAddr      = "ENV_NACOS_ADDR"
	EnvNacosPort      = "ENV_NACOS_PORT"
	EnvNacosNamespace = "ENV_NACOS_NAMESPACE"
	EnvNacosLogLevel  = "ENV_NACOS_LOG_LEVEL"
)

func getNacosConfig() ([]constant.ServerConfig, *constant.ClientConfig) {
	addrs := os.Getenv(EnvNacosAddr)
	if addrs == "" {
		addrs = "nacos-cs"
	}
	port := uint64(8848)
	if portStr := os.Getenv(EnvNacosPort); portStr != "" {
		_port, _ := strconv.ParseUint(portStr, 10, 64)
		if _port != 0 {
			port = uint64(_port)
		}
	}
	logLevel := os.Getenv(EnvNacosLogLevel)
	if logLevel == "" {
		logLevel = "debug"
	}

	sc := []constant.ServerConfig{
		*constant.NewServerConfig(addrs, port, constant.WithContextPath("/nacos")),
	}
	cc := constant.NewClientConfig(
		constant.WithNamespaceId(os.Getenv(EnvNacosNamespace)),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("/tmp/nacos/log"),
		constant.WithCacheDir("/tmp/nacos/cache"),
		constant.WithLogLevel(logLevel),
	)
	return sc, cc
}

func NewSourceNacos() (*SourceNacos, error) {
	sc, cc := getNacosConfig()
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "create nacos client error. config:%+v", cc)
	}
	return &SourceNacos{
		client: client,
	}, nil
}

func (source *SourceNacos) Name() string {
	return "nacos"
}

func (source *SourceNacos) DefaultGroup() string {
	return NacosDefaultGroup
}

func (source *SourceNacos) Get(group string, key string) (string, error) {
	return source.get(group, key)
}

func (source *SourceNacos) get(group string, key string) (string, error) {
	value, err := source.client.GetConfig(vo.ConfigParam{
		Group:  group,
		DataId: key,
	})
	if err != nil {
		return "", errors.Wrapf(err, "nacos get %s %s error.", group, key)
	}
	return value, nil
}

func (source *SourceNacos) BatchGet(group string, keys ...string) (map[string]string, error) {
	return source.batchGet(group, keys...)
}

func (source *SourceNacos) batchGet(group string, keys ...string) (map[string]string, error) {
	result := make(map[string]string, len(keys))
	var err error
	for _, key := range keys {
		result[key], err = source.get(key, group)
		if err != nil {
			return result, err
		}
	}
	return result, nil
}

func (source *SourceNacos) MonitorChange(group, key string, fn ConfigChangeCallback) error {
	if group == "" {
		group = NacosDefaultGroup
	}
	err := source.client.ListenConfig(vo.ConfigParam{
		Group:  group,
		DataId: key,
		OnChange: func(namespace, group, dataId, data string) {
			fn(data)
		},
	})
	if err != nil {
		err = errors.Wrapf(err, "nacos listen config error. key:%s, group:%s", key, group)
		return err
	}
	return nil
}

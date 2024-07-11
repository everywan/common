package configs

import (
	"context"
	"encoding/json"
	"strconv"
	"sync"

	"github.com/everywan/common/logger"
	"gopkg.in/yaml.v3"
)

var (
	defaultSource ISource
	defaultGroup  string

	// 初始化使用
	lock      = &sync.Mutex{}
	getSource func() ISource
)

func initSourceOnce() ISource {
	lock.Lock()
	defer lock.Unlock()
	if defaultSource == nil {
		initDefaultSource()
	}
	getSource = func() ISource {
		return defaultSource
	}
	return defaultSource
}

// 默认使用 nasos 实例. 如果使用 source_db 请通过 SetSource 实现
func initDefaultSource() {
	var err error
	defaultSource, err = NewSourceNacos()
	if err != nil {
		panic(err)
	}
}

func init() {
	getSource = initSourceOnce
}

func SetSource(source ISource) {
	defaultGroup = source.DefaultGroup()
	defaultSource = source
}

func Get(key string) string {
	return GetByGroup(defaultGroup, key)
}

func GetOrDefault(key string, defaultValue string) string {
	return GetByGroupOrDefault(defaultGroup, key, defaultValue)
}

func GetInt(key string) int {
	return GetIntByGroup(defaultGroup, key)
}

func GetIntOrDefault(key string, defaultValue int) int {
	return GetIntByGroupOrDefault(defaultGroup, key, defaultValue)
}

func GetJson(key string, out interface{}) {
	GetJsonByGroup(defaultGroup, key, out)
}

func GetYaml(key string, out interface{}) {
	GetYamlByGroup(defaultGroup, key, out)
}

func getByGroup(group, key string) (string, error) {
	value, err := getSource().Get(group, key)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. key:%s, err:%s", err, key)
		return "", err
	}
	return value, nil
}

func GetByGroup(group, key string) string {
	value, _ := getByGroup(group, key)
	return value
}

func GetByGroupOrDefault(group, key string, defaultValue string) string {
	value, err := getByGroup(group, key)
	if err != nil {
		return defaultValue
	}
	return value
}

func getIntByGroup(group, key string) (int, error) {
	value, err := getSource().Get(group, key)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. group:%s, key:%s, value:%s, err:%s",
			group, key, value, err)
		return 0, err
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. parse int error, group:%s, key:%s, value:%s, err:%s",
			group, key, value, err)
		return 0, err
	}
	return int(intValue), nil
}

func GetIntByGroup(group, key string) int {
	value, _ := getIntByGroup(group, key)
	return value
}

func GetIntByGroupOrDefault(group, key string, defaultValue int) int {
	value, err := getIntByGroup(group, key)
	if err != nil {
		value = defaultValue
	}
	return value
}

func GetJsonByGroup(group, key string, out interface{}) {
	value, err := getSource().Get(group, key)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. group:%s, key:%s, err:%s", group, key, err)
		return
	}
	err = json.Unmarshal([]byte(value), out)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. unmarshal error. group:%s, key:%s, value:%s, err:%s",
			group, key, value, err)
	}
}

func GetYamlByGroup(group, key string, out interface{}) {
	value, err := getSource().Get(group, key)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. group:%s, key:%s, err:%s", group, key, err)
		return
	}
	err = yaml.Unmarshal([]byte(value), out)
	if err != nil {
		logger.Errorf(context.TODO(), "get config error. unmarshal error. group:%s, key:%s, value:%s, err:%s",
			group, key, value, err)
	}
}

func MonitorChange(group, key string, callback ConfigChangeCallback) {
	getSource().MonitorChange(group, key, callback)
}

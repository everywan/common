package configs

type ConfigChangeCallback func(newvalue string)

//go:generate mockgen -destination mock/source.go -source=source.go
type ISource interface {
	Name() string
	DefaultGroup() string
	Get(group, key string) (string, error)
	BatchGet(group string, keys ...string) (map[string]string, error)
	MonitorChange(group, key string, fn ConfigChangeCallback) error
}

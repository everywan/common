package configs

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/everywan/common/logger"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

const (
	DBDefaultGroup = ""
)

var (
	tableName        = "prog_config"
	batchQueryDBSize = 100
)

type ProgConfig struct {
	ID    int64  `gorm:"column:id"`
	Group string `gorm:"column:group"`
	Name  string `gorm:"column:name"`
	Value string `gorm:"column:value"`
}

func (ProgConfig) TableName() string {
	return tableName
}

type SourceDB struct {
	db                 *gorm.DB
	sf                 *singleflight.Group
	monitorKeys        map[string]map[string]ConfigChangeCallback // map[group][key]
	localCacheGrowLock sync.Locker
	localCaches        map[string]*sync.Map // 分 group 获取时

	cronTickerDur time.Duration
}

var _ ISource = (*SourceDB)(nil)

func NewSourceDB(db *gorm.DB) *SourceDB {
	source := &SourceDB{
		db:                 db,
		sf:                 &singleflight.Group{},
		monitorKeys:        map[string]map[string]ConfigChangeCallback{},
		localCacheGrowLock: &sync.Mutex{},
		localCaches: map[string]*sync.Map{
			DBDefaultGroup: {}, // 初始化默认组
		},
		cronTickerDur: time.Second * 60 * 5,
	}
	go source.cronUpdate()
	return source
}

func (source *SourceDB) Name() string {
	return "db"
}

func (source *SourceDB) DefaultGroup() string {
	return DBDefaultGroup
}

func (source *SourceDB) growLocalCaches(group string) {
	source.localCacheGrowLock.Lock()
	defer source.localCacheGrowLock.Unlock()
	// 二次检测
	if source.localCaches[group] != nil {
		return
	}

	localCaches := make(map[string]*sync.Map, len(source.localCaches)+1)
	localCaches[group] = &sync.Map{}
	for group, cache := range source.localCaches {
		localCaches[group] = cache
	}
	source.localCaches = localCaches
}

func (source *SourceDB) Get(group, key string) (string, error) {
	if source.localCaches[group] == nil {
		source.growLocalCaches(group)
	}
	value, loaded := source.localCaches[group].Load(key)
	if loaded {
		return value.(string), nil
	}
	value, err, _ := source.sf.Do("get_"+key+"_"+group, func() (interface{}, error) {
		return source.get(group, key)
	})
	if err == nil {
		source.localCaches[group].Store(key, value)
	}
	return value.(string), err
}

func (source *SourceDB) get(group, key string) (string, error) {
	record := &ProgConfig{}
	query := source.db.Select("value").Where("group = ?", group)
	query = query.Where("name = ?", key)
	err := query.First(record).Error
	if err != nil {
		return "", errors.Wrapf(err, "query from source_db error. group:%s, key:%s", group, key)
	}
	return record.Value, nil
}

func (source *SourceDB) BatchGet(group string, keys ...string) (map[string]string, error) {
	if source.localCaches[group] == nil {
		source.growLocalCaches(group)
	}
	cache := source.localCaches[group]
	fetchKeys := []string{}
	result := make(map[string]string, len(keys))
	for _, key := range keys {
		value, loaded := cache.Load(key)
		if loaded {
			result[key] = value.(string)
		} else {
			fetchKeys = append(fetchKeys, key)
		}
	}
	fetchResult, err := source.batchGet(group, fetchKeys...)
	if err != nil {
		return result, err
	}
	for k, v := range fetchResult {
		result[k] = v
		cache.Store(k, v)
	}
	return result, nil
}

func (source *SourceDB) batchGet(group string, keys ...string) (map[string]string, error) {
	result := make(map[string]string, len(keys))
	for start := 0; start < len(keys); start += batchQueryDBSize {
		end := start + batchQueryDBSize
		if end > len(keys) {
			end = len(keys)
		}
		query := source.db.Model(&ProgConfig{}).Select("name,value")
		query = query.Where("group = ?", group)
		records := []ProgConfig{}
		err := query.Where("name in (?)", keys[start:end]).Find(&records).Error
		if err != nil {
			return result, errors.Wrapf(err, "update from source_db error, group:%s, keys:%v", group, keys[start:end])
		}
		for _, record := range records {
			result[record.Name] = record.Value
		}
	}
	return result, nil
}

// 非并发安全
func (source *SourceDB) MonitorChange(group, key string, fn ConfigChangeCallback) error {
	if _, ok := source.monitorKeys[group]; !ok {
		source.monitorKeys[group] = make(map[string]ConfigChangeCallback)
	}
	source.monitorKeys[group][key] = fn
	// 监听时先触发一遍获取
	_, err := source.Get(group, key)
	return err
}

func (source *SourceDB) cronUpdate() {
	randomDur := time.Duration(rand.Intn(10)) * time.Second
	ticker := time.NewTicker(source.cronTickerDur + randomDur)
	for range ticker.C {
		for group, localCache := range source.localCaches {
			source.updateConfigs(group, localCache)
		}
	}
}

func (source *SourceDB) updateConfigs(group string, localCache *sync.Map) {
	keys := []string{}
	localCache.Range(func(key, _ any) bool {
		keyStr := key.(string)
		keys = append(keys, keyStr)
		return true
	})

	kvs, err := source.batchGet(group, keys...)
	if err != nil {
		logger.Errorf(context.Background(), "updateConfigs error:%v, group:%s, keys:%+v",
			err, group, keys)
		// 宁可错过, 不可错误
		return
	}

	if source.monitorKeys[group] == nil {
		source.monitorKeys[group] = map[string]ConfigChangeCallback{}
	}

	for _, key := range keys {
		newValue := kvs[key]

		// 执行监听函数
		if fn, ok := source.monitorKeys[group][key]; ok {
			oldValue, _ := localCache.Load(key)
			if oldValue == nil {
				oldValue = ""
			}
			if oldValue != newValue {
				// 程序内部控制是否 panic
				go fn(newValue)
			}
		}

		// 存储新数据. 新数据删除时存空字符串
		localCache.Store(key, newValue)
	}
}

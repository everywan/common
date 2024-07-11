package configs

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/everywan/common/tests/mysql"
	"github.com/stretchr/testify/assert"
)

func TestSourceDBGet(t *testing.T) {
	suite := mysql.NewMysqlMock(t)
	defer suite.Close()

	source := NewSourceDB(suite.Gdb)
	tcases := []struct {
		name        string
		group, key  string
		expectValue string
	}{
		{"case1", "", "test1", "value1"},
		{"case2", DBDefaultGroup, "test2", "value2"},
		{"case3", "test_group", "test3", "value3"},
	}

	for _, tcase := range tcases {
		tcase := tcase
		if tcase.group == DBDefaultGroup {
			tcase.group = ""
		}
		t.Run(tcase.name, func(t *testing.T) {
			sql := fmt.Sprintf("SELECT `value` FROM `%s` WHERE group = \\? AND name = \\? ORDER BY `%s`.`id` LIMIT \\?", tableName, tableName)
			suite.Mock.ExpectQuery(sql).
				WithArgs(tcase.group, tcase.key, 1).
				WillReturnRows(
					sqlmock.NewRows([]string{"value"}).
						FromCSVString(tcase.expectValue),
				)
			value, err := source.Get(tcase.group, tcase.key)
			assert.NoError(t, err, "sourceDB.get error")
			assert.Equal(t, tcase.expectValue, value, "not expcet")

			value2, err := source.Get(tcase.group, tcase.key)
			assert.NoError(t, err, "get from cache error")
			assert.Equal(t, value, value2, "get from cache not equal")
		})
	}
}

func TestSourceDBMonitor(t *testing.T) {
	suite := mysql.NewMysqlMock(t)
	defer suite.Close()
	source := NewSourceDB(suite.Gdb)
	done := make(chan struct{})

	// 预设初始值
	expectValue := "value1"
	actualValue := expectValue
	testKey := "test"
	source.localCaches[DBDefaultGroup].Store(testKey, expectValue)
	source.MonitorChange(DBDefaultGroup, testKey, func(newvalue string) {
		actualValue = newvalue
		done <- struct{}{}
	})

	// 更改值
	go func() {
		expectValue = "value2"
		sql := fmt.Sprintf("SELECT name,value FROM `%s` WHERE group = \\? AND name in \\(\\?\\)", tableName)
		suite.Mock.ExpectQuery(sql).
			WithArgs(DBDefaultGroup, testKey).
			WillReturnRows(
				sqlmock.NewRows([]string{"name", "value"}).
					FromCSVString(testKey + "," + expectValue),
			)
		source.updateConfigs(DBDefaultGroup, source.localCaches[DBDefaultGroup])
	}()
	<-done
	assert.Equal(t, expectValue, actualValue)
}

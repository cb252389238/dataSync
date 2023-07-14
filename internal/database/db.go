package database

import (
	"fmt"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"riskDataSync/internal/config"
	"riskDataSync/util/tools"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once sync.Once
	sets *MysqlSets
)

type MysqlSets struct {
	mysql map[string]*gorm.DB
	l     sync.RWMutex
}

func (r *MysqlSets) Key(key ...string) *gorm.DB {
	r.l.RLock()
	defer r.l.RUnlock()
	var name string
	if len(key) <= 0 {
		name = "default"
	} else {
		name = key[0]
	}
	if db, ok := r.mysql[name]; ok {
		return db
	}
	return nil
}

func NewDb() *MysqlSets {
	once.Do(func() {
		conf := config.GetHotConf()
		path := tools.GetRootPath()
		dir := path + conf.LogPath
		_, err := tools.MakeDir(dir)
		if err != nil {
			panic(fmt.Sprintf("创建日志目录失败 path:%s", dir))
		}
		filename := strings.Replace(dir+"/"+"sql.log", "\\", "/", -1)
		fileIo, _ := os.Create(filename)
		//载入mysql集合
		dbClients := map[string]*gorm.DB{}
		for _, m := range conf.Mysql {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", m.User, m.Password, m.Host, m.Port, m.Database)
			slowLogger := logger.New(
				//将标准输出作为Writer
				log.New(fileIo, "\r\n", log.LstdFlags),
				logger.Config{
					//设定慢查询时间阈值为1ms
					SlowThreshold: 1 * time.Second,
					//设置日志级别，只有Warn和Info级别会输出慢查询日志
					LogLevel:                  logger.Warn,
					IgnoreRecordNotFoundError: true,
				},
			)
			open, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
				Logger: slowLogger,
			})
			if err != nil {
				panic(err)
			}
			db, err := open.DB()
			if err != nil {
				panic(err)
			}
			db.SetMaxOpenConns(5)
			db.SetMaxIdleConns(10)
			dbClients[m.Name] = open
		}
		sets = &MysqlSets{
			mysql: dbClients,
		}

	})
	return sets
}

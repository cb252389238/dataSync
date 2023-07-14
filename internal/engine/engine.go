package engine

import (
	"fmt"
	"os"
	"riskDataSync/internal/database"
	"riskDataSync/internal/offset"
	"riskDataSync/util/cache"
	"riskDataSync/util/log"
	"riskDataSync/util/task"
	"sync"
	"time"
)

type Engine struct {
	Wg      *sync.WaitGroup
	L       *sync.RWMutex
	Mysql   *database.MysqlSets
	Offsets *offset.Offsets
	Cache   *cache.Cache
	Log     *log.LocalLogger
	Tasks   *task.Task
}

func NewEngine() *Engine {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	ctx := &Engine{
		Wg:      &sync.WaitGroup{},
		L:       &sync.RWMutex{},
		Mysql:   database.NewDb(),
		Offsets: offset.New(),
		Cache:   cache.New(cache.NoExpiration, time.Second*300),
		Log:     log.NewLog(),
		Tasks:   task.New(),
	}
	for _, w := range ctx.Tasks.GetAllTask() {
		for _, b := range w.Battle {
			ctx.Offsets.Set(w.Name+b.Name, b.StartVal)
		}
	}
	return ctx
}

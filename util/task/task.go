package task

import (
	"encoding/json"
	"fmt"
	"os"
	"riskDataSync/util/tools"
	"sync"
)

var (
	t *Task
)

type War struct {
	Source_db string   `json:"source_db"`
	Target_db string   `json:"target_db"`
	Battle    []Battle `json:"task"`
	Name      string   `json:"name"`
}

type Battle struct {
	Name               string `json:"name"`
	Source_table       string `json:"source_table"`
	Target_table       string `json:"target_table"`
	Fields             string `json:"fields"`
	Source_primary_key string `json:"source_primary_key"`
	Target_primary_key string `json:"target_primary_key"`
	OffsetField        string `json:"offset_field"`
	StartVal           int    `json:"start_val"`
	EndVal             int    `json:"end_val"`
	Where              string `json:"where"`
}

type Task struct {
	Data    []War
	l       sync.RWMutex
	taskNum int
	offset  int
}

func New() *Task {
	jsonPath := tools.GetRootPath() + "/task.json"
	file, err := os.Open(jsonPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var data = []War{}
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		panic(err)
	}
	for _, w := range data {
		for _, b := range w.Battle {
			if b.StartVal < 0 || b.EndVal < 0 {
				panic(fmt.Sprintf("配置错误：%s=>%s start_val或end_val不能小于0", w.Name, b.Name))
			}
			if b.EndVal < b.StartVal && b.EndVal != 0 {
				panic(fmt.Sprintf("配置错误：%s=>%s end_val不能小于start_val", w.Name, b.Name))
			}
		}
	}
	t = &Task{
		Data:    data,
		taskNum: len(data),
	}
	return t
}

func (t *Task) GetAllTask() []War {
	t.l.RLock()
	defer t.l.RUnlock()
	return t.Data
}

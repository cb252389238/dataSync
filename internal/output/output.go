package output

import (
	"fmt"
	"riskDataSync/internal/engine"
	"riskDataSync/util/tools"
	"time"
)

func Print(ectx *engine.Engine) {
	t := time.NewTicker(time.Millisecond * 200)
	tasks := ectx.Tasks.GetAllTask()
	for {
		select {
		case <-t.C:
			tools.Cls()
			for _, v := range tasks {
				fmt.Printf("任务名：%s\r\n", v.Name)
				for _, t := range v.Battle {
					ectx.Offsets.L.RLock()
					offsetNum := ectx.Offsets.Get(v.Name + t.Name)
					ectx.Offsets.L.RUnlock()
					fmt.Printf("【%s】偏移位置：%d\r\n", t.Name, offsetNum)
				}
			}
		}
	}
}

func EndPrint(ectx *engine.Engine) {
	fmt.Println("-----------------程序运行结束-----------------")
	tasks := ectx.Tasks.GetAllTask()
	for _, v := range tasks {
		fmt.Printf("任务名：%s\r\n", v.Name)
		for _, t := range v.Battle {
			ectx.Offsets.L.RLock()
			offsetNum := ectx.Offsets.Get(v.Name + t.Name)
			ectx.Offsets.L.RUnlock()
			fmt.Printf("【%s】偏移位置：%d\r\n", t.Name, offsetNum)
		}
	}
}

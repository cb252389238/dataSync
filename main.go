package main

import (
	"fmt"
	"os"
	"riskDataSync/internal/config"
	"riskDataSync/internal/engine"
	"riskDataSync/internal/logic"
	"riskDataSync/internal/output"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()
	config.LoadConfig()        //载入配置文件
	ectx := engine.NewEngine() //载入引擎
	go output.Print(ectx)
	//开始同步数据
	ectx.Wg.Add(1)
	go logic.Start(ectx)
	ectx.Wg.Wait()
	output.EndPrint(ectx)
}

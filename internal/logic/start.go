package logic

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"riskDataSync/internal/engine"
	"riskDataSync/typedef"
	"riskDataSync/util/task"
	"strings"
	"time"
)

func recoverErr(ectx *engine.Engine) {
	if err := recover(); err != nil {
		ectx.Log.Error("panic:%+v", err)
	}
}

func Start(ectx *engine.Engine) {
	defer ectx.Wg.Done()
	tasks := ectx.Tasks.GetAllTask()
	for _, v := range tasks {
		ectx.Wg.Add(1)
		go war(ectx, v)
	}
}

func war(ectx *engine.Engine, war task.War) {
	defer ectx.Wg.Done()
	for _, v := range war.Battle {
		ectx.Wg.Add(1)
		go battle(ectx, v, war)
	}
}

func battle(ectx *engine.Engine, battle task.Battle, war task.War) {
	defer ectx.Wg.Done()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	fields := parseFields(battle.Fields)
	fields.sourceFields = append(fields.sourceFields, battle.OffsetField)
	for {
		select {
		case <-ctx.Done():
			return
		case typedef.ChCoroutine <- 1:
			ectx.Wg.Add(1)
			go execute(ectx, battle, war, cancel, fields)
		}
	}
}

func execute(ectx *engine.Engine, battle task.Battle, war task.War, cancelFunc context.CancelFunc, fields FieldsStruct) {
	defer func() {
		ectx.Wg.Done()
		<-typedef.ChCoroutine
		recoverErr(ectx)
	}()
	key := war.Name + battle.Name
	ectx.Offsets.L.Lock()
	offsetNum := ectx.Offsets.Get(key)
	where := fmt.Sprintf("%s > %d", battle.OffsetField, offsetNum)
	if battle.Where != "" {
		where = fmt.Sprintf("%s and %s", where, battle.Where)
	}
	if battle.EndVal > 0 {
		where = fmt.Sprintf("%s and %s <= %d", where, battle.OffsetField, battle.EndVal)
	}
	data := []map[string]any{}
	err := ectx.Mysql.Key(war.Source_db).Debug().Table(battle.Source_table).Select(fields.sourceFields).Where(where).Order(battle.OffsetField + " asc").Limit(500).Find(&data).Error
	if err != nil {
		ectx.Offsets.L.Unlock()
		ectx.Log.Error("任务:%s=>%s查询出错:%+v", war.Name, battle.Name, err)
		return
	}
	if len(data) > 0 {
		ectx.Offsets.Set(key, cast.ToInt(data[len(data)-1][battle.OffsetField]))
		special(battle.Name, battle.Source_table, data)
	} else {
		ectx.Offsets.L.Unlock()
		cancelFunc()
		return
	}
	ectx.Offsets.L.Unlock()
	targetDatas := []map[string]any{}
	for _, d := range data {
		target := map[string]any{}
		for _, field := range fields.res {
			if v, ok := d[field.sname]; ok {
				if field.stype == field.ttype {
					target[field.tname] = v
				} else {
					switch field.stype {
					case "int": //整型类型
						switch field.ttype {
						case "int": //整型类型
							target[field.tname] = v
						case "float": //浮点类型
							target[field.tname] = cast.ToFloat64(v)
						case "string":
							target[field.tname] = cast.ToString(v)
						case "year": //yyyy 例如2023
							target[field.tname] = cast.ToTime(v).Format("2006")
						case "time": //HH:MM:SS
							target[field.tname] = cast.ToTime(v).Format("15:04:05")
						case "date": //YY-MM-DD
							target[field.tname] = cast.ToTime(v).Format("20060102")
						case "datetime": //YY-MM-DD HH:MM:SS
							target[field.tname] = cast.ToTime(v).Format("2006-01-02 15:04:05")
						case "timestamp":
							target[field.tname] = cast.ToTime(data).Format("2006-01-02 15:04:05")
						}
					case "float": //浮点类型
						switch field.ttype {
						case "int": //整型类型
							target[field.tname] = cast.ToInt(v)
						case "float": //浮点类型
							target[field.tname] = v
						case "string":
							target[field.tname] = cast.ToString(v)
						}
					case "string":
						switch field.ttype {
						case "int": //整型类型
							target[field.tname] = cast.ToInt(v)
						case "float": //浮点类型
							target[field.tname] = cast.ToFloat64(v)
						case "string":
							target[field.tname] = v
						}
					case "year": //yyyy 例如2023
						switch field.ttype {
						case "string":
							target[field.tname] = cast.ToString(v)
						case "year": //yyyy 例如2023
							target[field.tname] = v
						}
					case "time": //HH:MM:SS
						switch field.ttype {
						case "string":
							target[field.tname] = cast.ToString(v)
						case "time":
							target[field.tname] = v
						}
					case "date": //YY-MM-DD
						switch field.ttype {
						case "string":
							target[field.tname] = cast.ToString(v)
						case "date":
							target[field.tname] = v
						}
					case "datetime": //YY-MM-DD HH:MM:SS
						layout := "2006-01-02 15:04:05 -0700 MST" // 时间字符串格式
						t, _ := time.Parse(layout, cast.ToString(v))
						switch field.ttype {
						case "int": //整型类型
							target[field.tname] = cast.ToTime(t).Unix()
						case "string":
							target[field.tname] = cast.ToString(v)
						case "year": //yyyy 例如2023
							target[field.tname] = cast.ToTime(v).Format("2006")
						case "datetime": //YY-MM-DD HH:MM:SS
							target[field.tname] = v
						}
					case "timestamp":
						layout := "2006-01-02 15:04:05" // 时间字符串格式
						t, _ := time.Parse(layout, cast.ToString(v))
						switch field.ttype {
						case "int": //整型类型
							target[field.tname] = cast.ToTime(t).Unix()
						case "string":
							target[field.tname] = cast.ToString(v)
						case "year": //yyyy 例如2023
							target[field.tname] = cast.ToTime(t).Format("2006")
						case "time": //HH:MM:SS
							target[field.tname] = cast.ToTime(t).Format("15:04:05")
						case "date": //YY-MM-DD
							target[field.tname] = cast.ToTime(t).Format("20060102")
						case "datetime": //YY-MM-DD HH:MM:SS
							target[field.tname] = cast.ToTime(t).Format("2006-01-02 15:04:05")
						case "timestamp":
							target[field.tname] = cast.ToTime(t).Format("2006-01-02 15:04:05")
						}

					}
				}
			}
		}
		targetDatas = append(targetDatas, target)
	}
	err = ectx.Mysql.Key(war.Target_db).Debug().Table(battle.Target_table).Create(&targetDatas).Error
	if err != nil {
		ectx.Log.Error("任务:%s=>%s插入出错:%+v", war.Name, battle.Name, err)
	}
}

type FieldsStruct struct {
	source       [][]string
	sourceFields []string
	target       [][]string
	targetFields []string
	res          []FieldRes
}

type FieldRes struct {
	sname string
	tname string
	stype string
	ttype string
}

func parseFields(str string) FieldsStruct {
	fields := FieldsStruct{}
	split := strings.Split(str, "|")
	for _, a1 := range split {
		a2 := strings.Split(a1, "=>")
		if len(a2) != 2 {
			panic(fmt.Sprintf("任务配置错误:%s", str))
		}
		source := strings.Split(a2[0], ",")
		target := strings.Split(a2[1], ",")
		fields.target = append(fields.target, target)
		fields.source = append(fields.source, source)
		fields.sourceFields = append(fields.sourceFields, source[0])
		fields.targetFields = append(fields.targetFields, target[0])
		fieldRes := FieldRes{
			sname: source[0],
			tname: target[0],
			stype: source[1],
			ttype: target[1],
		}
		fields.res = append(fields.res, fieldRes)
	}
	return fields
}

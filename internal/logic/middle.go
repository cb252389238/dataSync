package logic

import "github.com/spf13/cast"

func special(taskName, sourceTableName string, data []map[string]any) {
	for k, v := range data {
		for k1, v1 := range v {
			data[k][k1] = rule(taskName, sourceTableName, k1, v1)
		}
	}
}

func rule(taskName, tableName, fieldName string, fieldVal any) any {
	switch taskName {
	case "kh_appraise_order_case_log":
		switch fieldName {
		case "case_res":
			return cast.ToInt(fieldVal) + 1
		}
	case "kh_pw_appraise_praise1":
		switch fieldName {
		case "is_useful":
			return 1
		}
	case "kh_pw_appraise_praise2":
		switch fieldName {
		case "is_praise":
			return 2
		}
	}
	return fieldVal
}

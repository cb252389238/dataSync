# dataSync
* 数据同步

 1. 首先配置config.yaml，在项目根目录下。

```yaml

LogPath: logs #日志目录
Mysql:  #mysql的配置
  -
    Host: 127.0.0.1 #注意冒号后面都有空格
    Port: 3307
    User: root
    Password: root
    Database: test
    Name: default #别名，注意别名后面在任务配置的时候会用到，代表的是连接那个数据库
  -
    Host: 127.0.0.1
    Port: 3307
    User: root
    Password: root
    Database: test2
    Name: risk

```

2. 配置task.json  配置要同步的任务。在项目根目录下

```json5
[
  {
    "name": "task-1", //主任务，一般涉及多个数据库同步会有多个主任务
    "source_db": "default", //源数据库别名，和config.yaml里的别名对应
    "target_db": "risk", //目标数据库别名
    "task": [ //子任务 两个数据库内不同table同步的任务
      {
        "name": "task-1-1", //子任务名
        "source_table": "dict1",  //源表名
        "target_table": "dict3",  //目标表名
        "fields": "word,string=>name,string", //需要同步的字段格式下面详细说明
        "source_primary_key": "id", //源表主键，没有可留空
        "target_primary_key": "id", //目标主键，没有可留空
        "offset_field": "id", //偏移字段，依靠该字段作为游标和排序的索引
        "start_val": 0, //偏移字段开始值，不能小于0，不包含初始值从初始值+1位置开始
        "end_val": 0, //偏移字段结束值。为0会同步到最后一条数据。如果不为0不能小于start_val。取值范围<end_val，不包含本身
        "where": "status =0 and type = 0" //源表查询条件
      },
      {
        "name": "task-1-2",
        "source_table": "dict2",
        "target_table": "dict4",
        "fields": "word1,string=>name1,string|word2,string=>name2,string",
        "source_primary_key": "id",
        "target_primary_key": "id",
        "offset_field": "id",
        "start_val": 10,
        "end_val": 0,
        "where": ""
      }
    ]
  }
]
```

4. fields字段配置说明

```
"fields": "word,string=>name,string"|"word,string=>name,string",
源表字段名,源表字段类型=>目标字段名,目标字段类型
多个字段映射用|隔开
```
**支持转换的类型**

| 源类型    | 目标类型   |
|:-------|:-------|
| string | int    |
| date   | string |
| time   | string |
| time   | int    |
| date   | int    |
| timestamp    | string |
| timestamp    | int    |


5. 将编译好的文件随意放入任何位置，然后将config.yaml、task.json放到可执行程序根目录下。执行可执行程序即可


package typedef

import "sync"

type HotConf struct {
	Conf Config
	L    sync.RWMutex
}

type Config struct {
	Mysql   []Mysql
	LogPath string
}

type Mysql struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Name     string
}

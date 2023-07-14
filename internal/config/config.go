package config

import (
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"riskDataSync/typedef"
	"runtime"
)

var (
	hotConf    typedef.HotConf
	configPath string
)

func init() {
	configPath = getRootPath()
}

func LoadConfig() {
	var conf typedef.Config
	configFile := "config"
	configFileExt := "yaml"
	c := viper.New()
	c.SetConfigName(configFile)
	c.SetConfigType(configFileExt)
	c.AddConfigPath(configPath)
	err := c.ReadInConfig()
	if err != nil {
		panic("配置文件载入错误" + configPath + configFile + "." + configFileExt)
	}
	err = c.Unmarshal(&conf)
	if err != nil {
		panic("配置文件解析错误" + err.Error())
	}
	hotConf.L.Lock()
	hotConf.Conf = conf
	hotConf.L.Unlock()
}

func getAbsBinPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	rst := filepath.Dir(path)
	return rst + "/"
}

func getUserBinPath() string {
	var cpath string = ""
	path, _ := os.Getwd()
	cpath = path + "/"
	return cpath
}

func getRootPath() string {
	if runtime.GOOS == "linux" {
		return getAbsBinPath()
	} else {
		return getUserBinPath()
	}
}

func GetHotConf() typedef.Config {
	hotConf.L.RLock()
	defer hotConf.L.RUnlock()
	conf := hotConf.Conf
	return conf
}

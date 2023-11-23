package config

import (
	"os"
	"sync"
	"time"

	rlog "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	config GlobalConfig //全局业务配置文件
	once   sync.Once
)

// db配置结构

type DbConf struct {
	Host          string `yaml:"host" mapstructure:"host"`
	Port          string `yaml:"port" mapstructure:"port"`
	User          string `yaml:"user" mapstructure:"user"`
	Password      string `yaml:"password" mapstructure:"password"`
	Dbname        string `yaml:"dbname" mapstructure:"dbname"`
	Max_idle_conn int    `yaml:"max_idle_conn" mapstructure:"max_idle_conn"`
	Max_open_conn int    `yaml:"max_open_conn" mapstructure:"max_open_conn"`
	Max_idle_time int64  `yaml:"max_idle_time" mapstructure:"max_idle_time"`
}

// type DbConf struct {
// 	Host          string `yaml:"host"`
// 	Port          string `yaml:"port"`
// 	User          string `yaml:"user"`
// 	Password      string `yaml:"password"`
// 	Dbname        string `yaml:"dbname"`
// 	Max_idle_conn int    `yaml:"max_idle_conn"`
// 	Max_open_conn int    `yaml:"max_open_conn"`
// 	Max_idle_time int64  `yaml:"max_idle_time"`
// }

// log配置结构
type LogConf struct {
	LogPattern string `yaml:"log_pattern" mapstructure:"log_pattern"`
	LogPath    string `yaml:"log_path" mapstructure:"log_path"`
	SaveDays   uint   `yaml:"save_days"mapstructure:"save_days"`
	Level      string `yaml:"level" mapstructure:"level"`
}

// AppConf 服务配置
type AppConf struct {
	AppName string `yaml:"app_name" mapstructure:"app_name"` //业务名
	Version string `yaml:"version" mapstructure:"version"`   //版本
	Port    int    `yaml:"port" mapstructure:"port"`         //端口
	RunMode string `yaml:"run_mode" mapstructure:"run_mode"` //运行模式
}

type GlobalConfig struct {
	DbConfig  DbConf  `yaml:"db" mapstructure:"db"`   //db配置
	LogConfig LogConf `yaml:"log" mapstructure:"log"` //log配置
	AppConf   AppConf `yaml:"app" mapstructure:"app"`
}

func readConf() {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./conf")
	viper.AddConfigPath("../conf/")
	err := viper.ReadInConfig()
	if err != nil {
		panic("read  config file err:" + err.Error())
	}
	err = viper.Unmarshal(&config)
	if err != nil {
		panic("config file unmarshal err:" + err.Error())
	}
	// fmt.Println(config)
}

func GetGlobalConf() *GlobalConfig {
	once.Do(readConf)
	return &config
}

// 日志初始化
func InitConfig() {
	globalConf := GetGlobalConf()
	//设置日志级别
	level, err := log.ParseLevel(globalConf.LogConfig.Level)
	if err != nil {
		panic("log level parseerr:" + err.Error())
	}
	//设置日志格式为json格式
	log.SetFormatter(&logFormatter{
		log.TextFormatter{
			DisableColors:   true,                  //禁止颜色输出
			TimestampFormat: "2006-01-02 15:04:05", //设置时间戳格斯
			FullTimestamp:   true,                  //时间戳包含日期和时间
		}})
	log.SetReportCaller(true) //用于设置是否在日志中报告调用者的信息。如果将其设置为 true，日志中将包含调用者的文件名和行号等信息
	log.SetLevel(level)
	switch globalConf.LogConfig.LogPattern {
	case "stdout":
		log.SetOutput(os.Stdout)
	case "stderr":
		log.SetOutput(os.Stderr)
	case "file":
		if _, err := os.Stat(globalConf.LogConfig.LogPath); os.IsNotExist(err) {
			//日志目录不存在，创建目录
			err := os.MkdirAll(globalConf.LogConfig.LogPath, 0775)
			if err != nil {
				panic("failed to create log dirctory:" + err.Error())
			}
		}
		logger, err := rlog.New(
			globalConf.LogConfig.LogPath+".%Y%m%d%H%M", //日志名
			// rlog.WithRotationCount(globalConf.LogConfig.SaveDays),
			rlog.WithMaxAge(24*time.Hour*time.Duration(globalConf.LogConfig.SaveDays)), //日志保留的最大天数
			rlog.WithRotationTime(time.Hour*24),                                        //日志切割的时间
		)
		if err != nil {
			panic("log conf err:" + err.Error())
		}
		log.SetOutput(logger)
	default:
		panic("log conf err,checklog_pattern in log.yaml")
	}
}

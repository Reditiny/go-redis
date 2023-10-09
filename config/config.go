package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type config struct {
	Server *ServerConfig `mapstructure:"server" json:"server"`
	Logger *LogConfig    `mapstructure:"log" json:"log"`
}

// Conf 全局配置变量
var Conf = new(config)

// InitConfig 设置读取配置信息
func InitConfig() {
	// 获取当前工作目录
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}
	// 初始化 viper 单例模式 指定配置文件
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/")
	// 读取配置信息
	err = viper.ReadInConfig()
	// 保存配置信息到配置类
	saveConf(Conf)
}

func saveConf(Conf *config) {
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Conf); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}
	fmt.Println("保存配置成功")
}

type ServerConfig struct {
	Bind string `mapstructure:"bind" json:"bind"`
	Port int    `mapstructure:"port" json:"port"`
}

type LogConfig struct {
	Level int `mapstructure:"level" json:"level"`
}

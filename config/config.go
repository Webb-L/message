package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type ServiceConfig struct {
	App struct {
		Title    string `yaml:"title"`
		Debug    bool   `yaml:"debug"`
		Language string `yaml:"language"`
		Store    struct {
			Path   string `yaml:"path"`
			Prefix string `yaml:"prefix"`
		} `yaml:"store"`
		Env string `yaml:"env"`
		Log struct {
			Info   string `yaml:"info"`
			Error  string `yaml:"error"`
			Access string `yaml:"access"`
		} `yaml:"log"`
	} `yaml:"app"`
	Database struct {
		Host       string `yaml:"host"`
		Port       int    `yaml:"port"`
		User       string `yaml:"user"`
		Pwd        string `yaml:"pwd"`
		Name       string `yaml:"name"`
		MaxIdleCon int    `yaml:"max_idle_con"`
		MaxOpenCon int    `yaml:"max_open_con"`
		Params     struct {
			Character string `yaml:"character"`
		} `yaml:"params"`
	} `yaml:"database"`
	API struct {
		Test     bool `yaml:"test"`
		MaxLimit int  `yaml:"maxLimit"`
	} `yaml:"api"`
}

var AppConfig ServiceConfig

func InitConfig() {
	workDir, _ := os.Getwd()
	// 设置配置文件的名字
	viper.SetConfigName("config")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件的路径，指定 config 目录下寻找
	viper.AddConfigPath(workDir + "/config")
	// 寻找配置文件并读取
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		panic(fmt.Errorf("fatal error parset file: %w", err))
	}
}

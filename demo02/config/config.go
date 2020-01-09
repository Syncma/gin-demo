package config

import (
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
}

func Init(cfg string) error {
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name) // 如果指定了配置文件，则解析指定的配置文件
	} else {
		viper.AddConfigPath("conf") // 如果没有指定配置文件，则解析默认的配置文件
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")               // 设置配置文件格式为YAML
	viper.AutomaticEnv()                      // 读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER")           // 读取环境变量的前缀为APISERVER
	replacer := strings.NewReplacer(".", "_") //这里面是做什么用的？
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper解析配置文件
		return err
	}

	return nil
}

// 监控配置文件变化并热加载程序
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("Config file changed: %s", e.Name)
	})
}

/*
这里发现修改了配置文件config.yaml里面的端口号，程序并没有使用新的端口？

自动监听配置是这个意思：如果修改了配置，下次再用viper.GetString之类的函数，获取配置时，获取到的就是最新的值。
改了配置端口后，程序其实没用重新去Get新的值并注册新的端口，所以端口还是旧的


*/

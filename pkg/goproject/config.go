package goproject

import (
	"fmt"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// SourceControl configuration sets the default upstream initializing git projects
type SourceControl struct {
	Uri string
}

// Go sets various Go Modules related options
type Go struct {
	Vendor bool
}

// Template specifices user template locations
type Template struct {
	Name string `mapstructure:"name"`
	Path string `mapstructure:"path"`
}

// Config application configuration
type Config struct {
	SourceControl   SourceControl `mapstructure:"sourcecontrol"`
	Go              Go            `mapstructure:"go"`
	CustomTemplates []Template    `mapstructure:"custom_templates"`
}

// Load the configuration from disk
func Load() (*Config, error) {
	var c Config
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(path.Join(home, ".config", "goproject"))
	viper.SetConfigName("config")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&c)
	return &c, err
}

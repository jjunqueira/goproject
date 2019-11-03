package goproject

import (
	"fmt"
	"os"
)

type GoProject struct {
	Config *Config
}

func New() GoProject {
	c, err := Load()
	if err != nil {
		fmt.Println("Unable to load configuration file %v", err)
		os.Exit(-1)
	}
	return GoProject{Config: c}
}

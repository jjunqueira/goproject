package goproject

import (
	"fmt"
	"os"
)

// GoProject contains base application dependencies
type GoProject struct {
	Config *Config
}

// New constructs a new instance of the application.
func New() GoProject {
	c, err := Load()
	if err != nil {
		fmt.Printf("Unable to load configuration file %v", err)
		os.Exit(-1)
	}

	return GoProject{Config: c}
}

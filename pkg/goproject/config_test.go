package goproject

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	c, err := Load()
	if err != nil {
		t.Errorf("%v", err)
	}
	fmt.Printf("%v", c)
}

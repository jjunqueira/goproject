package goproject

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	_, err := Load()
	if err != nil {
		t.Errorf("%v", err)
	}
}

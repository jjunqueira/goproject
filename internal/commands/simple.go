package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const simpleMain = `
package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`

// Simple Generates a simple project template
func Simple(gitPrefix string, projectname string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	fullPath := filepath.Join(cwd, projectname)
	err = os.Mkdir(fullPath, 0777)
	if err != nil {
		return err
	}

	err = gitInit(fullPath)
	if err != nil {
		return err
	}

	err = initGoModule(fullPath, gitPrefix, projectname)
	if err != nil {
		return err
	}

	err = writeMain(filepath.Join(fullPath, "main.go"), simpleMain)
	return err
}

func gitInit(dir string) error {
	cmd := exec.Command("git", "init", "-q", dir)
	return cmd.Run()
}

func initGoModule(dir string, gitPrefix string, projectname string) error {
	moduleName := fmt.Sprintf("%s/%s", gitPrefix, projectname)
	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir
	return cmd.Run()
}

func writeMain(path string, contents string) error {
	bytes := []byte(contents)
	err := ioutil.WriteFile(path, bytes, 0644)
	return err
}

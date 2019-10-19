package templates

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const emptyMain string = `
package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`

// Empty generates a simple project template
func Empty(gitPrefix string, projectname string) error {
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

	err = writeMain(filepath.Join(fullPath, "main.go"), emptyMain)
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
	tmpPath := path + ".tmp"
	err := ioutil.WriteFile(tmpPath, bytes, 0644)
	if err != nil {
		return err
	}

	cmd := exec.Command("gofmt", tmpPath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to get stdout pipe for gofmt %v", err))
	}

	err = cmd.Start()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to start gofmt %v", err))
	}

	b, err := ioutil.ReadAll(stdout)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to read formatted main.go file %v", err))
	}

	err = ioutil.WriteFile(path, b, 0660)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to write formatted main.go %v", err))
	}

	err = cmd.Wait()
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to finish gofmt command %v", err))
	}

	err = os.Remove(tmpPath)
	if err != nil {
		return errors.New(fmt.Sprintf("Unable to delete temporary main file %v", err))
	}

	return nil
}

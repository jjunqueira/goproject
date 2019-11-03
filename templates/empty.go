/*
Copyright © 2019 Joshua Junqueira <joshua.junqueira@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package templates

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

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

	err = copyTemplate("empty", fullPath)
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

func copyTemplate(tplName string, dest string) error {
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	tplPath := path.Join(home, ".config", "goproject", "templates", tplName)
	if err != nil {
		return err
	}
	err = copyFiles(tplPath, dest)
	return err
}

func copyFiles(src string, dest string) error {
	cmd := exec.Command("cp", "-r", path.Join(src, "*"), path.Join(dest, "*"))
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

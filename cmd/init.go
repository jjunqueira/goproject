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
package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const configTpl string = `[sourcecontrol]
uri = "https://github.com"

[go]
vendor = false

[[custom_templates]]
name = "example1"
path = "/path/to/custom/templates/example1"

[[custom_templates]]
name = "example2"
path = "/path/to/custom/templates/example2"
`

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize goproject for the first time",
	Long: `The init command initializes goproject to run for the first time. By default it will do the following:
1. Create a directory in the users home folder .config/goproject
2. Create a default config .config/goproject/config.toml
3. Create a templates directory .config/goproject/templates
4. Download default templates to .config/goproject/templates

example: goproject init
	`,
	Run: func(cmd *cobra.Command, args []string) {
		err := initialize()
		if err != nil {
			fmt.Printf("Unable to initialize goproject: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func initialize() error {
	var err error

	var configPath string

	defer func() {
		if err != nil && configPath != "" {
			err = os.RemoveAll(configPath)
			if err != nil {
				fmt.Printf("unable to clean up files %v", err)
			}
		}
	}()

	configPath, err = defaultPath()
	if err != nil {
		return fmt.Errorf("unable to construct configuration path: %v", err)
	}

	fmt.Printf("Initializing goproject to %s\n", configPath)

	stat, err := os.Stat(configPath)
	if err == nil && stat.IsDir() {
		errMsg := `the provided path already exists so this command will do nothing.
		If you want to start over from scratch remove the configuration directory and rerun this command
		`

		return errors.New(errMsg)
	}

	// Create the main configuration directory
	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("unable to create the configuration directory %v", err)
		return err
	}

	fmt.Printf("Creating default configuration file %s\n", path.Join(configPath, "config.toml"))

	// Write default config
	bytes := []byte(configTpl)
	tmpPath := path.Join(configPath, "config.toml")

	err = ioutil.WriteFile(tmpPath, bytes, 0644)
	if err != nil {
		err = fmt.Errorf("unable to create the configuration file %v", err)
		return err
	}

	fmt.Printf("Creating templates directory %s\n", path.Join(configPath, "templates"))

	err = os.MkdirAll(path.Join(configPath, "templates"), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("unable to create the templates directory %v", err)
		return err
	}

	fmt.Printf("Downloading base templates to %s\n", path.Join(configPath, "templates"))

	err = downloadTemplates(path.Join(configPath, "templates"))
	if err != nil {
		err = fmt.Errorf("unable to download templates %v", err)
		return err
	}

	return nil
}

func defaultPath() (string, error) {
	home, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	return path.Join(home, ".config", "goproject"), nil
}

func downloadTemplates(tplPath string) error {
	cmd := exec.Command("git", "clone", "https://github.com/jjunqueira/goproject-templates", tplPath)
	return cmd.Run()
}

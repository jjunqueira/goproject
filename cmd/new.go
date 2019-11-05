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
	"fmt"
	"os"

	"github.com/jjunqueira/goproject/pkg/goproject"
	"github.com/jjunqueira/goproject/pkg/templates"
	"github.com/spf13/cobra"
)

var gitPrefix string

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [template] [projectname]",
	Short: "Create a new project based on a template",
	Long:  `Create a new project based on one of the available templates such as 'empty' or 'cli'`,
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app = goproject.New()
		p, err := templates.NewProject(app.Config, gitPrefix, args[0], args[1])
		if err != nil {
			fmt.Printf("Unable to construct project structure: %v", err)
			os.Exit(1)
		}

		err = templates.Generate(app.Config, p)
		if err != nil {
			fmt.Printf("Unable to generate project files: %v", err)
			os.Exit(2)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.PersistentFlags().StringVarP(&gitPrefix, "gitprefix", "", "", "The git prefix to use for the project e.g. github.com/jjunqueira. This will be prepended to the project name for use with Go modules. If no prefix is provided the project name will be used as the name of the module.")
}

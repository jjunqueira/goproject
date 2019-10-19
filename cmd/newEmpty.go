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

	"github.com/jjunqueira/goproject/templates"

	"github.com/spf13/cobra"
)

// newEmptyCmd represents the newEmpty command
var newEmptyCmd = &cobra.Command{
	Use:   "empty [projectname]",
	Short: "Creates a project based on the 'empty' template",
	Long: `Creates a project based on the 'empty' template.

The empty template is the most basic project template.
It creates the root directory, initializes git, initializes Go module, and adds a main.go file`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a project name argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		templates.Empty(args[0], args[0])
	},
}

func init() {
	newCmd.AddCommand(newEmptyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	newEmptyCmd.PersistentFlags().String("gitprefix", "", "The git prefix to use for the project e.g. github.com/jjunqueira. This will be prepended to the project name for use with Go modules. If no prefix is provided the project name will be used as the name of the module.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newEmptyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

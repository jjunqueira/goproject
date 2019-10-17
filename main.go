package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/jjunqueira/goproject/internal/commands"
)

const usage string = "Usage: goproject new TEMPATETYPE PROJECTNAME"

func main() {

	if len(os.Args) < 2 {
		println(usage)
		os.Exit(-1)
	}

	commandType := os.Args[1]
	var err error
	switch commandType {
	case "new":
		err = new(os.Args)
	case "tidy":
		println("'tidy' not yet supported")
		os.Exit(-1)
	}

	if err != nil {
		fmt.Printf("%v\n%s", err, usage)
	}
}

func new(inputArgs []string) error {
	if len(inputArgs) < 3 {
		return errors.New("Not enough input arguments")
	}

	command := inputArgs[2]

	switch command {
	case "simple":
		if len(inputArgs) < 4 {
			return errors.New("Not enough input arguments")
		}
		projectName := inputArgs[3]
		return commands.Simple("github.com/jjunqueira", projectName)
	default:
		return errors.New("Project type not supported")
	}
}

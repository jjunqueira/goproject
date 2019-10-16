package main

import "github.com/jjunqueira/goproject/internal/commands"

func main() {
	commands.Echo("vim-go")
	if true {
		commands.Echo("true")
	}
}

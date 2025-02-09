package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading input: ", err)
			os.Exit(1)
		}
		commandHandler(command[:len(command)-1])
	}
}

func commandHandler(command string) {
	cmds := strings.Split(command, " ")
	if len(cmds) == 0 {
		return
	}
	switch cmds[0] {
	case "exit":
		exitCommand()
	case "echo":
		echoCommand(cmds[1:])
	default:
		invalidCommand(command)
	}
}

func echoCommand(cmds []string) {
	s := strings.Join(cmds, " ")
	fmt.Println(s)
}

func invalidCommand(command string) {
	fmt.Println(command + ": command not found")
}

func exitCommand() {
	os.Exit(0)
}

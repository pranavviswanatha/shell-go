package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var handlerMap map[string]func([]string)

func main() {
	initMap()
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
	f, ok := handlerMap[cmds[0]]
	if ok {
		f(cmds[1:])
		return
	}
	cmd := exec.Command(cmds[0], cmds[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		invalidCommand(cmds)
	}

}

func echoCommand(cmds []string) {
	s := strings.Join(cmds, " ")
	fmt.Println(s)
}

func invalidCommand(cmds []string) {
	fmt.Fprintln(os.Stdout, cmds[0]+": command not found")
}

func exitCommand(cmds []string) {
	os.Exit(0)
}

func typeCommand(cmds []string) {
	_, ok := handlerMap[cmds[0]]
	if ok {
		fmt.Fprintln(os.Stdout, cmds[0]+" is a shell builtin")
		return
	}
	paths := strings.Split(os.Getenv("PATH"), ":")
	for _, path := range paths {
		fp := filepath.Join(path, cmds[0])
		if _, err := os.Stat(fp); err == nil {
			fmt.Println(fp)
			return
		}
	}
	fmt.Fprintln(os.Stdout, cmds[0]+": not found")
}

func pwdCommand(cmds []string) {
	path, _ := os.Getwd()
	fmt.Fprintln(os.Stdout, path)
}

func cdCommand(cmds []string) {
	err := os.Chdir(cmds[0])
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s: No such file or directory\n", cmds[0])
	}
}

func initMap() {
	handlerMap = make(map[string]func([]string))
	handlerMap["exit"] = exitCommand
	handlerMap["echo"] = echoCommand
	handlerMap["type"] = typeCommand
	handlerMap["pwd"] = pwdCommand
	handlerMap["cd"] = cdCommand
}

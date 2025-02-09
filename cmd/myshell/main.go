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

func splitCommand(command string) []string {
	var cmds []string
	s := strings.Trim(command, "\r\n")
	for {
		start := strings.IndexAny(s, "'\"")
		if start == -1 {
			cmds = append(cmds, strings.Fields(s)...)
			break
		}
		// fmt.Println("log: pointcrossed")
		ch := s[start]
		cmds = append(cmds, strings.Fields(s[:start])...)
		s = s[start+1:]
		end := strings.IndexByte(s, ch)
		token := s[:end]
		if len(s) > end+1 && s[end+1] == ' ' {
			token = token + " "
		}
		cmds = append(cmds, token)
		s = s[end+1:]
	}
	return cmds
}

func commandHandler(command string) {
	cmds := splitCommand(command)
	if len(cmds) == 0 {
		return
	}
	f, ok := handlerMap[cmds[0]]
	if ok {
		go f(cmds[1:])
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
	s := strings.Join(cmds, "")
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
	if len(cmds) != 1 {
		return
	}
	path := cmds[0]
	if strings.TrimSpace(path) == "~" {
		path = os.Getenv("HOME")
	}
	err := os.Chdir(path)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%s: No such file or directory\n", path)
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

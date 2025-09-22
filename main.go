package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Alias struct {
	Name    string
	Command string
}

type EnvVar struct {
	Key string
	Val string
}

type Config struct {
	Aliases []Alias
	EnvVars []EnvVar
	Path    string
}

func (C *Config) Parse(path string) {
	C.Path = path
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("error opening config file", err)
		os.Exit(1)
	}
	defer file.Close()
	x, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("error reading config file", err)
		os.Exit(1)
	}
	data := string(x)
	C.Aliases = make([]Alias, 0, 50)
	for line := range strings.Lines(data) {
		if strings.HasPrefix(line, "alias") {
			aliasLine := strings.TrimPrefix(line, "alias ")
			name, command, _ := strings.Cut(aliasLine, "=")
			command = strings.TrimSpace(command)
			C.Aliases = append(C.Aliases, Alias{name, command})
		}
	}
}

func (C *Config) AddLine(line string) {
	data, err := os.ReadFile(C.Path)
	if err != nil {
		fmt.Println("error reading config file", err)
		os.Exit(1)
	}
	strData := string(data)
	strData += "\n" + line
	os.WriteFile(C.Path, []byte(strData), 0644)
}

func (C *Config) PrintAliases() {
	data, err := os.ReadFile(C.Path)
	if err != nil {
		fmt.Println("error reading config file", err)
		os.Exit(1)
	}
	strData := string(data)
	cnt := 1
	for line := range strings.Lines(strData) {
		if aliasLine, found := strings.CutPrefix(line, "alias "); found == true {
			fmt.Printf("%d. %s\n", cnt, strings.TrimSpace(aliasLine))
			cnt++
		}
	}

}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("usage :\n\t- zrc add <line_to_add_to_shell_config>\n\t- zrc list aliases")
		os.Exit(0)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir", err)
		os.Exit(1)
	}

	path := filepath.Join(home, "Downloads/zsc.txt")
	var config Config
	config.Parse(path)
	var line string
	switch args[1] {
	case "add":
		for i, arg := range args {
			if i > 1 {
				line += arg + " "
			}
		}
		config.AddLine(line)
	case "list":
		if args[2] == "aliases" {
			config.PrintAliases()
		}
	}
}

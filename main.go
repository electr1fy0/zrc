package main

import (
	"fmt"
	"os"
	"os/exec"
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
	Aliases  []Alias
	EnvVars  []EnvVar
	Path     string
	Contents string
}

func (C *Config) Parse(path string) {
	C.Path = path

	x, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("error reading config file:", err)
		os.Exit(1)
	}

	data := string(x)
	C.Contents = data
}

func (C *Config) AddLine(line string) {
	C.Contents += "\n" + line
	if err := os.WriteFile(C.Path, []byte(C.Contents), 0644); err != nil {
		fmt.Println("error writing file:", err)
		os.Exit(1)
	}
}

func (C *Config) AddToPath(path string) {
	line := "export PATH=$PATH:" + path
	C.AddLine(line)
}

func (C *Config) AddAlias(key, value string) {
	line := fmt.Sprintf(`alias %s="%s"`, key, value)
	C.AddLine(line)
}

func (C *Config) RemoveAlias(key string) {
	for line := range strings.Lines(C.Contents) {
		if aliasLine, found := strings.CutPrefix(line, "alias "); found {
			name, _, _ := strings.Cut(aliasLine, "=")
			name = strings.TrimSpace(name)
			if name == key {
				C.RemoveLine(line)
				return
			}
		}
	}
}

func (C *Config) RemoveFromPath(path string) {
	formatted := "export PATH=$PATH:" + path
	for line := range strings.Lines(C.Contents) {
		if strings.TrimSpace(line) == formatted {
			C.RemoveLine(line)
		}
	}
}

func (C *Config) RemoveLine(toRemove string) {
	lines := strings.Split(C.Contents, "\n")
	var newContents []string
	for _, line := range lines {
		if strings.TrimSpace(line) != strings.TrimSpace(toRemove) {
			newContents = append(newContents, line)
		}
	}
	C.Contents = strings.Join(newContents, "\n")
	if err := os.WriteFile(C.Path, []byte(C.Contents), 0644); err != nil {
		fmt.Println("error writing file:", err)
		os.Exit(1)
	}
}

func (C *Config) PrintAliases() {
	data, err := os.ReadFile(C.Path)
	if err != nil {
		fmt.Println("error reading config file:", err)
		os.Exit(1)
	}
	strData := string(data)
	cnt := 1
	fmt.Println("Aliases:")
	for line := range strings.Lines(strData) {
		if aliasLine, found := strings.CutPrefix(line, "alias "); found {
			fmt.Printf("%d. %s\n", cnt, strings.TrimSpace(aliasLine))
			cnt++
		}
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage:")
		fmt.Println("\tzrc add <line>                  - add a raw line to shell config")
		fmt.Println("\tzrc alias add <name> <command>  - add an alias")
		fmt.Println("\tzrc alias remove <name>         - remove an alias")
		fmt.Println("\tzrc path add <dir>              - add directory to PATH")
		fmt.Println("\tzrc path remove <dir>           - remove directory from PATH")
		fmt.Println("\tzrc list aliases                - list all aliases")
		os.Exit(0)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("error getting home dir:", err)
		os.Exit(1)
	}

	appConfigPath := filepath.Join(home, ".zrcc")
	var shell string
	data, err := os.ReadFile(appConfigPath)
	if err != nil {
		fmt.Print("enter your shell config file name (e.g. .zshrc): ")
		fmt.Scanln(&shell)
		if err := os.WriteFile(appConfigPath, []byte(shell), 0644); err != nil {
			fmt.Println("error writing .zrcc:", err)
			os.Exit(1)
		}
	} else {
		shell = strings.TrimSpace(string(data))
	}

	path := filepath.Join(home, shell)
	var config Config
	config.Parse(path)

	cmd := exec.Command("source ", path)

	switch args[1] {
	case "add":
		if len(args) < 3 {
			fmt.Println("usage: zrc add <line_to_add>")
			os.Exit(0)
		}
		line := strings.Join(args[2:], " ")
		config.AddLine(line)
		err = cmd.Run()

	case "alias":
		if len(args) < 4 {
			fmt.Println("usage:")
			fmt.Println("\tzrc alias add <name> <command>")
			fmt.Println("\tzrc alias remove <name>")
			os.Exit(0)
		}
		switch args[2] {
		case "add":
			name := args[3]
			command := strings.Join(args[4:], " ")
			config.AddAlias(name, command)
			err = cmd.Run()
		case "remove":
			name := args[3]
			config.RemoveAlias(name)
			err = cmd.Run()
		default:
			fmt.Println("unknown alias command:", args[2])
		}

	case "path":
		if len(args) < 4 {
			fmt.Println("usage:")
			fmt.Println("\tzrc path add <dir>")
			fmt.Println("\tzrc path remove <dir>")
			os.Exit(0)
		}
		switch args[2] {
		case "add":
			dir := args[3]
			config.AddToPath(dir)
			err = cmd.Run()
		case "remove":
			dir := args[3]
			config.RemoveFromPath(dir)
			err = cmd.Run()
		default:
			fmt.Println("unknown path command:", args[2])
		}

	case "list":
		if args[2] == "aliases" {
			config.PrintAliases()
		} else {
			fmt.Println("unknown list command:", args[2])
		}

	default:
		fmt.Println("unknown command:", args[1])
	}
	if err != nil {
		fmt.Println("Error refreshing terminal based on config")
	}

}

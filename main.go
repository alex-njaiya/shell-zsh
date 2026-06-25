package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// preffix
const colorPreffix = "\033["
const (
	coloReset   string = colorPreffix + "0m"
	colorGreen  string = colorPreffix + "32m"
	colorBlue   string = colorPreffix + "34m"
	colorCyan   string = colorPreffix + "36m"
	colorYellow string = colorPreffix + "33m"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	// listen for every write from the keyboard
	for {
		path, err := getpath()

		if err != nil {
			fmt.Printf("%s > %s", colorGreen, coloReset)
			continue
		}

		if path == "/" {
			fmt.Printf("%s/ > %s", colorGreen, coloReset)
		} else {
			fmt.Printf("%s%s/ %s> %s", colorGreen, path, colorYellow, coloReset)
		}

		//read the keyboard input
		input, err := reader.ReadString('\n')

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// handle the input execution
		if err = execInput(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func execInput(input string) error {
	// remove the new line characte at the end of the input
	args := strings.Fields(input)

	if len(args) == 0 {
		return nil
	}

	command := args[0]
	arguments := args[1:]

	// check for built in commands
	switch command {
	case "cd":
		// check the length of the input
		// if the length is less than 2 get the root directory as the fallback
		// check the length of args if less than 2 throw a path error
		if len(arguments) == 0 {
			// if the length of argument is one. Go to the homedir
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			// if no error change to the homedir and update the command prompt
			return os.Chdir(homeDir)
		}

		// if there is an argument change to that specific argument
		return os.Chdir(arguments[0])
	case "exit", "Exit":
		os.Exit(0)
	}

	// prepare the command to be executed
	cmd := exec.Command(command, arguments...)

	// set the correct output device
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	// execute the command and return the error
	return cmd.Run()
}

func getpath() (path string, err error) {
	cdir, err := os.Getwd()

	if err != nil {
		return "> ", err
	}

	// format the homeDir path to use a tilde instead of the entire path
	path, err = formatHomeDirPath(cdir)
	return path, err
}

func formatHomeDirPath(target string) (path string, err error) {
	// get the homedirpath
	fullPath, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	if strings.HasPrefix(target, fullPath) {
		// get the base
		tilde := strings.Replace(target, fullPath, "~", 1)
		return tilde, err
	}

	return target, nil

}

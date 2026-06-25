package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)


func main() {
	reader := bufio.NewReader(os.Stdin)

	// listen for every write from the keyboard
	for {
		path, err := getpath()
		if err != nil {
			fmt.Print("> ")
		} 
		fmt.Print(path, "/ > ")
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
			homeDir, err := os.UserHomeDir()

			if err != nil {
				if err := os.Chdir("/"); err != nil {
						return err
					}
					return err
			}
			return os.Chdir(homeDir)
		}

		// change the dir using os.Chdir
		// whenever I change a directory I update the new line format to include its name
		return os.Chdir(arguments[0])
	case "exit", "Exit":
		os.Exit(0)
	}

	// prepare the command to be executed
	cmd := exec.Command(command, arguments...)

	// set the correct output device
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	// execute the command and return the error
	return cmd.Run()
}


func getpath() (path string, err error) {
	cdir, err := os.Getwd() 

	if err != nil {
		return "> ", err
	}

	// extract folder name
	dirName := filepath.Base(cdir)

	return dirName, err
}
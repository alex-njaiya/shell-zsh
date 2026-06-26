package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
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

var (
	history      []string
	historyIndex int
	currentInput string
)

func main() {

	// fd for the standard input
	fd := int(os.Stdin.Fd())

	// put the terminal into raw mode
	oldState, err := term.MakeRaw(fd)

	if err != nil {
		log.Fatal(err)
	}

	// restore the terminal state when main exists or finishes reading
	defer term.Restore(fd, oldState)

	buf := make([]byte, 3)
	// listen for every write from the keyboard
outer:
	for {
		currentInput = ""
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
	inner:
		for {
			n, err := readInput(buf)

			if err != nil {
				break outer
			}

			// check for specific escape sequences
			if n == 3 && buf[0] == 27 && buf[1] == '[' {
				if buf[2] == 'A' {
					// up arrow pressed. Cycle to previous command history
					// // decrement the history index
					if historyIndex > 0 {
						historyIndex -= 1

					}

					// // render history[historyIndex]
					currentInput = history[historyIndex]
					fmt.Print(currentInput)
					continue
				}

				if buf[2] == 'B' {
					// Down arrow pressed. Cycle to the next command history
					continue
				}
			}

			//check for the enterkey to execute command.
			// In raw mode enter sends a \r carriage return or newline
			if buf[0] == '\r' || buf[0] == '\n' {
				fmt.Print("\n")

				// save to history and update the histryIndex
				if currentInput != "" {
					history = append(history, currentInput)
					historyIndex = len(history)
				}
				break inner // break out of the reading loop to execute the command
			}

			if buf[0] == '\x7f' {
				// delete the last char
				if len(currentInput) > 0 {
					// remove the last character from the internal string tracker
					currentInput = currentInput[:len(currentInput)-1]

					// move the cursor back, overwrite with space move cursor back
					fmt.Print("\b \b")
				}
				continue
			}

			currentInput += string(buf[:n])
			fmt.Print(string(buf[:n]))

		}

		// handle the input execution
		if err = execInput(currentInput); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func readInput(buf []byte) (input int, err error) {

	n, err := os.Stdin.Read(buf)

	if err != nil {
		return 0, err
	}

	return n, nil
}

func execInput(input string) error {
	// remove the new line characte at the end of the input
	// TODO: We have to redo this because we are
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

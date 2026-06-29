package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/term"
	"text-based-shell/utils"
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
	history         []string
	historyIndex    int
	currentInput    string
	unexecutedInput string
)

var filename = "command-history.txt"

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

	// listen for every write from the keyboard
outer:
	for {
		buf := make([]byte, 3)

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
					if len(history) > 0 && historyIndex > 0 {

						if historyIndex == len(history) {
							unexecutedInput = currentInput
						}
						historyIndex -= 1
						currentInput = history[historyIndex]

						// 1.\r moves cursor to start of line
						// 2. \033[K clears everything from cursor position to the end of the line
						fmt.Print("\r\033[K")

						// Reprint your prompt first so it doesn't disappear
						rePrintPrompt(path)
						fmt.Print(currentInput)
					}
					continue
				}

				if buf[2] == 'B' {
					// Down arrow pressed. Cycle to the next command history
					if historyIndex < len(history) {
						historyIndex += 1
						// clear any input when the down button is pressed
						fmt.Print("\r\033[K")

						rePrintPrompt(path)

						if historyIndex == len(history) {
							// if the hisrory index == end of the history replace with the current input text
							// keep track of the current input sent so as to display after the end of the history
							currentInput = unexecutedInput
						} else {
							currentInput = history[historyIndex]
						}
						fmt.Print(currentInput)

					}
					continue
				}
			}

			//check for the enterkey to execute command.
			// In raw mode enter sends a \r carriage return or newline
			if buf[0] == '\r' || buf[0] == '\n' {
				fmt.Print("\r\n")

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
				} else {
					// if the current input is empty do nothing
				}
				continue
			}

			if buf[0] == 3 {
				fmt.Print("\r\n")
				break outer
			}

			if buf[0] >= 32 && buf[0] != 127 {
				currentInput += string(buf[:n])
				fmt.Print(string(buf[:n]))
				continue
			}

		}

		// temprarily restore the terminal to normal mode
		term.Restore(fd, oldState)

		// handle the input execution
		if err = execInput(currentInput); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		// CLEAR: Reset variables here. Clear drafts only after a command finishes executing
		currentInput = ""
		unexecutedInput = ""

		// re-enable raw mode immediately so the shell can read keys
		oldState, err = term.MakeRaw(fd)

		if err != nil {
			log.Fatal(err)
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
		err := saveHistory(utils.Write{
			Filename: filename,
			History:  history,
		})

		if err != nil {
			fmt.Printf("Error saving configuration history: %v", err)
		} else {
			fmt.Println("Command session saved successfully")
		}

		history = nil
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

func rePrintPrompt(path string) {
	if path == "/" {
		fmt.Printf("\r%s/ > %s", colorGreen, coloReset)
	} else {
		fmt.Printf("\r%s%s/ %s> %s", colorGreen, path, colorYellow, coloReset)
	}
}

func saveHistory(writer utils.FileWriter) error {
	if writer == nil {
		return fmt.Errorf("Cannot save history: Provided writer is nil")
	}

	if err := writer.WriteToFile(); err != nil {
		return fmt.Errorf("Failed to save history: %v", err)
	}

	return nil
}

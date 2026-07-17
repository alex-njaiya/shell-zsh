package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	// "syscall"

	"text-based-shell/utils"

	"golang.org/x/term"
)

// preffix
const colorPreffix = "\033["
const (
	coloReset  string = colorPreffix + "0m"
	colorGreen string = colorPreffix + "32m"
	colorBlue   string = colorPreffix + "34m"
	// colorCyan   string = colorPreffix + "36m"
	colorYellow string = colorPreffix + "33m"
)

var (
	history         []string
	sessionHistory  []string
	historyIndex    int
	currentInput    string
	unexecutedInput string
)

const (
	backspace     = '\x7f'
	escapeChar    = 27
	escapeBracket = '['
	maxBufSize    = 3
)

var cursor int

var filename = "command-history.txt"
var foregroundCmd *exec.Cmd
var path string

func main() {
	// fd for the standard input
	fd := int(os.Stdin.Fd())


	// run setup 
	cfg := utils.RunSetup(fd)



	inputChan := make(chan []byte)
	interruptChan := make(chan struct{}, 1)
	var err error

	go func() {
		for {
			buf := make([]byte, 3)
			n, err := os.Stdin.Read(buf)

			if err != nil {
				close(inputChan)
				return
			}
			inputChan <- buf[:n]
		}
	}()

	// load the history file into memory
	logger := &utils.Write{
		Filename: filename,
		History:  history,
	}

	history, err = logger.ReadFromFile()

	if err != nil {
		fmt.Printf("Error loading history: %v\n", err)
	}

	// put the terminal into raw mode
	oldState, err := term.MakeRaw(fd)

	if err != nil {
		log.Fatal(err)
	}

	// restore the terminal state when main exists or finishes reading
	defer term.Restore(fd, oldState)

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, os.Interrupt)

	go func() {
		for s := range sig {
			if foregroundCmd != nil {
				forePgid := foregroundCmd.Process.Pid

				syscall.Kill(-forePgid, s.(syscall.Signal))
			} else {
				fmt.Print("^C\r\n")
				currentInput = ""
				cursor = 0
				select {
				case interruptChan <- struct{}{}:
				default:
				}

				os.Stdin.Write([]byte{0})

			}
		}
	}()

	// listen for every write from the keyboard

outer:
	for {
		path, err := getpath()

		// append the username and host
		userspace := utils.AppendUsername(cfg)

		

		if err != nil {
			fmt.Printf("%s%s$", colorGreen, coloReset)
			continue
		}

		if path == "/" {
			fmt.Printf("%s%s%s%s$ %s", colorGreen, userspace, colorBlue, colorYellow, coloReset)
		} else {
			fmt.Printf("%s%s%s%s%s$ %s", colorGreen, userspace, colorBlue, path, colorYellow, coloReset)
		}

		//read the keyboard input
		// input channel
	inner:
		for {
			var buf []byte
			var n int

			select {
			case b, ok := <-inputChan:
				if !ok {
					break outer
				}

				select {
				case <-interruptChan:
					break inner
				default:
				}

				buf = b
				n = len(b)
			case <-interruptChan:
				currentInput = ""
				cursor = 0
				break inner
			}

			// check for specific escape sequences
			if n == maxBufSize && buf[0] == escapeChar && buf[1] == escapeBracket {
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
						rePrintPrompt(path, userspace)
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

						rePrintPrompt(path, userspace)

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

				// moving left with the cursor
				if buf[2] == 'D' {
					// check if the cursor it at zero
					// if at zero block it from going further

					if cursor > 0 {
						// decrement the cursor index by 1
						cursor--
						// move the terminal cursor left
						fmt.Print("\033[1D")
					} else {
						fmt.Print("\007")
					}

				}

				if buf[2] == 'C' {
					if cursor < len(currentInput) {
						cursor++
						// move the terminal cursor right
						fmt.Print("\033[1C")
					} else {
						fmt.Print("\007")
					}
				}
			}

			//check for the enterkey to execute command.
			// In raw mode enter sends a \r carriage return or newline
			if buf[0] == '\r' || buf[0] == '\n' {
				fmt.Print("\r\n")

				// save to history and update the histryIndex
				if currentInput != "" {
					history = append(history, currentInput)
					sessionHistory = append(sessionHistory, currentInput)
					historyIndex = len(history)
					cursor = 0
				}
				break inner // break out of the reading loop to execute the command
			}

			if buf[0] == '\x7f' {
				// delete the last char
				if len(currentInput) > 0 && cursor > 0 {
					// remove the last character from the internal string tracker
					currentInput = currentInput[:cursor-1] + currentInput[cursor:]
					cursor--
					// redraw the line
					fmt.Print("\r\033[K")
					rePrintPrompt(path, userspace)
					fmt.Print(currentInput)

					// move the terminal cursor back to the correct position
					correctCursorPos := len(currentInput) - cursor

					if correctCursorPos > 0 {
						fmt.Printf("\033[%dD", correctCursorPos)

					}
				} else {
					// if the current input is empty do nothing
				}
				continue
			}

			// tab space -- autocomplete
			if buf[0] == '\t' {
				if len(currentInput) == 0 {
					continue
				}
				// split the string into tokens
				tokens := strings.Fields(currentInput)

				if len(tokens) == 0 {
					// no tokens -- no keyboard input
					continue
				}
				// identify prefix and completion types

				isCommand := len(tokens) == 1 && !strings.HasSuffix(currentInput, " ")

				if isCommand {
					matches := utils.CompletionFromPath(tokens[0])
					//if the completion
					if len(matches) == 0 {
						// produce a beep sound
						fmt.Print("\007")
					}

					if len(matches) == 1 {
						currentInput = matches[0]
						// redraw the line
						fmt.Print("\r\033[K")
						rePrintPrompt(path, userspace)
						fmt.Print(currentInput)
					}

					if len(matches) > 1 {
						longestCommonPrefix := utils.LongestCommonPrefix(matches)
						currentInput = longestCommonPrefix
						fmt.Print("\r\n")

						fmt.Print(strings.Join(matches, " "))

						fmt.Print("\r\n")
						rePrintPrompt(path, userspace)
						fmt.Print(currentInput)
					}
				} else {
					// longest common prefix
					completions := utils.CompletionsFromFiles(tokens[len(tokens)-1])

					if len(completions) == 0 {
						fmt.Print("\a")
					}

					if len(completions) == 1 {
						base := strings.Join(tokens[:len(tokens)-1], " ")

						if base != "" {
							currentInput = base + " " + completions[0]
						} else {
							currentInput = completions[0]
						}
						fmt.Print("\r\033[K")
						rePrintPrompt(path, userspace)
						fmt.Print(currentInput)
					}

					if len(completions) > 1 {
						longestCommonPrefix := utils.LongestCommonPrefix(completions)
						base := strings.Join(tokens[:len(tokens)-1], " ")

						if base != "" {
							currentInput = base + " " + longestCommonPrefix
						} else {
							currentInput = longestCommonPrefix
						}
						fmt.Print("\r\n")
						fmt.Print(strings.Join(completions, " "))
						fmt.Print("\r\n")
						rePrintPrompt(path, userspace)
						fmt.Print(currentInput)
					}
				}

			}

			if buf[0] >= 32 && buf[0] != 127 {
				currentInput = currentInput[:cursor] + string(buf[:n]) + currentInput[cursor:]
				cursor++
				fmt.Print("\r\033[K")
				rePrintPrompt(path, userspace)
				fmt.Print(currentInput)
				// move the terminal back to the correct position
				correctCursorPos := len(currentInput) - cursor

				if correctCursorPos > 0 {
					fmt.Printf("\033[%dD", correctCursorPos)

				}
				continue
			}
		}

		// temporarily restore the terminal to normal mode
		term.Restore(fd, oldState)

		// handle the input execution
		if err = execInput(currentInput); err != nil {
			fmt.Fprintf(os.Stderr, "\r%s\r\n", err)
		}

		// CLEAR: Reset variables here. Clear drafts only after a command finishes executing
		currentInput = ""
		cursor = 0
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
		logger := &utils.Write{
			Filename: filename,
			History:  history,
		}

		err := logger.WriteToFile()

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

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	foregroundCmd = cmd

	if err := cmd.Start(); err != nil {
		return err
	}

	// wait for the process to be started
	if err := cmd.Wait(); err != nil {
		if _, ok := err.(*exec.ExitError); !ok {
			return err
		}
	}
	foregroundCmd = nil
	return nil
}

func getpath() (path string, err error) {
	cdir, err := os.Getwd()

	if err != nil {
		return "$ ", err
	}

	// format the homeDir path to use a tilde instead of the entire path
	path, err = formatHomeDirPath(cdir)

	// append the username and host on the path
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

func rePrintPrompt(path string, userspace string) {
	if path == "/" {
		fmt.Printf("\r%s%s%s%s$ %s", colorGreen, userspace, colorBlue, colorYellow, coloReset)
	} else {
		fmt.Printf("\r%s%s%s%s%s$ %s", colorGreen, userspace, colorBlue, path, colorYellow, coloReset)
	}
}

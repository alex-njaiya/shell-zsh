package utils

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

type Config struct {
	username string
	Password string
}

// it is supposed to read and parse the file
func LoadConfig() (Config, error) {
	var cfg Config
	dirPath, err := ConfigPath()

	if err != nil {
		return cfg, fmt.Errorf("Error getting the home dir path: %w", err)
	}

	configFilePath := filepath.Join(dirPath, "settings.conf")

	// read the file contents
	read, err := os.ReadFile(configFilePath)

	if err != nil {
		return cfg, fmt.Errorf("Error while reading the file contents: %w", err)
	}
	// split by newlines to get each line
	lines := strings.Split(string(read), "\n")

	// For each line split on = to get the key value
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		fmt.Print(parts)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			fmt.Printf("DEBUG -- Key: %s | value: %s", key, value)

			// assign to the right field on the struct based on the key

			switch key {
			case "username":
				cfg.username = value
			case "password":
				cfg.Password = value
			}
		}
	}

	// return the populated struct
	return cfg, nil
}

func SaveConfig(cfg Config) error {
	// take the config and save it to disk
	dirPath, err := ConfigPath()

	if err != nil {
		return fmt.Errorf("Error getting home directory path: %w", err)
	}

	configFilePath := filepath.Join(dirPath, "settings.conf")

	// format the config username=value\npassword=value\n
	formattedConfig := fmt.Sprintf("username=%s\npassword=%s\n", cfg.username, cfg.Password)

	// write the config to file using os.WriteFile
	// 0644 -- read and write permission for owner
	if err = os.WriteFile(configFilePath, []byte(formattedConfig), 0644); err != nil {
		return fmt.Errorf("Error writing the config file to disk: %w", err)
	}

	fmt.Println("Successfully saved configuration to disk")
	return nil
}

func SetUpExists() bool {
	path, err := ConfigPath()

	if err != nil {
		return false
	}

	_, err = os.Stat(path)

	if err != nil {
		return false
	}

	return true
}

func RunSetup(fd int) Config {
	var cfg Config
	// welcome message
	welcomeMsg := `
====================================================================
			WELCOME TO THE GOSH SHELL
====================================================================
Version: 1.0.0
Status: Connected Successfully

Please enter your desired username and password
to get started in using the gosh shell
---------------------------------------------------------------------
	`

	fmt.Print(welcomeMsg)

	// Prompt for a username and read it
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\r\033[K")
	fmt.Print("username: ")

	// capture the username
	input, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Error reading user input: %s", err)
		return cfg
	}

	// trim the the new line character at the end
	username := strings.TrimSpace(input)

	// Prompt for the user password using term.ReadPassword
	fmt.Print("password: ")
	passwordBytes, err := term.ReadPassword(fd)

	if err != nil {
		fmt.Printf("Error reading user password input: %s\n", err)
		return cfg
	}
	fmt.Println()  // move cursor to a new line after hidden inputs

	// Prompt to confirm the password and verify it matches the first
	fmt.Print("confirm password: ")
	confirmPasswordBytes, err := term.ReadPassword(fd)

	if err != nil {
		fmt.Printf("Error reading user confirm Password Input: %s\n", err)
		return cfg
	}
	// verify the password. Compare raw byte slices before hashing
	if string(passwordBytes) != string(confirmPasswordBytes) {
		fmt.Println("Passwords do not match. Please restart setup")
		return cfg
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(confirmPasswordBytes ,bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("Error hashing password: %s\n", err)
		return cfg
	}

	// build and populate the config struct
	cfg.username = username
	cfg.Password = string(hashedPassword)

	// save the config file onto disk
	err = SaveConfig(cfg)

	if err != nil {
		fmt.Printf("Error saving the config file into disk: %s\n", err)
	}

	return cfg
}

func VerifyPassword(cfg Config, input string) bool {
	return cfg.Password == strings.TrimSpace(input)
}

func ConfigPath() (string, error) {
	// get the home directory and append the /.gosh/config
	home, err := os.UserHomeDir()

	if err != nil {
		return "", err
	}

	targetDir := filepath.Join(home, ".gosh", "config")

	err = os.MkdirAll(targetDir, 0755)

	if err != nil {
		return "", err
	}

	return targetDir, nil
}

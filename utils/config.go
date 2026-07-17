package utils

import (
	"os"
	"path/filepath"
)

type Config struct {
	username string
	Password string
}

//it is supposed to read and parse the file
func LoadConfig() {

}

func SaveConfig() {

}

func SetUpExists(){

}

func RunSetup(){

}

func VerifyPassword(cfg Config, input string) {
	
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
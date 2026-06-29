package utils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

func IsfileEmpty(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, fmt.Errorf("File does not exist: %w", err)
		}

		return false, err
	}

	// check is the file is zero bytes
	return fileInfo.Size() == 0, nil
}

func AppendToFile(filename string, lines []string) error {
	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND

	file, err := os.OpenFile(filename, flags, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	writer := bufio.NewWriter(file)


	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		
		if err != nil  {
			return fmt.Errorf("Failed to write line to buffer")
		}
	}

	// push everything from memory to disk before exiting
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Failed to flush buffer to disk: %w", err)
	}

	return nil
}

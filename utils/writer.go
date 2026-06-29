package utils

import (
	"bufio"
	"fmt"
	"os"
)

type FileWriter interface {
	WriteToFile() error
}

type Write struct {
	Filename string
	History  []string
}

func IsfileEmpty(filename string) (bool, error) {
	fileInfo, err := os.Stat(filename)

	if os.IsNotExist(err) {
		return true, nil
	}

	if err != nil {
		return false, err
	}

	// check is the file is zero bytes
	return fileInfo.Size() == 0, nil
}

func (w Write) WriteToFile() error {
	if len(w.History) == 0 {
		return nil
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND

	// if empty create the file
	file, err := os.OpenFile(w.Filename,flags, 0644)

	if err != nil {
		return fmt.Errorf("Error creating file: %w", err)
	}

	defer file.Close()

	// if file is empty write to it
	writer := bufio.NewWriter(file)
	var totalBytes int

	for _, line := range w.History {
		n, err := writer.WriteString(line + "\n")

		if err != nil {
			return fmt.Errorf("Error writing to file: %v", err)
		}
		totalBytes += n
	}

	// push everything from memory to disk before exiting
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("Failed to flush buffer to disk: %w", err)
	}

	fmt.Printf("Successfully written %d bytes to %s\n", totalBytes, w.Filename)
	return nil
}
package utils

import (
	"bufio"
	"fmt"
	"os"
)

type HistoryManager interface {
	WriteToFile() error
	ReadFromFile() ([]string, error)
}

type Write struct {
	Filename string
	History  []string
}

func (w Write) WriteToFile() error {
	if len(w.History) == 0 {
		return nil
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_APPEND

	// if empty create the file
	file, err := os.OpenFile(w.Filename, flags, 0644)

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

func (w Write) ReadFromFile() ([]string, error) {
	//open the file for reading
	file, err := os.Open(w.Filename)

	if os.IsNotExist(err) {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("Failed to open history file: %w", err)
	}

	defer file.Close()

	reader := bufio.NewScanner(file)

	// reset the internal history slice
	w.History = []string{}

	for reader.Scan() {
		w.History = append(w.History, reader.Text())
	}

	// check whether the scanner encountered any errors during reading
	if err := reader.Err(); err != nil {
		return nil, fmt.Errorf("Error reading history file lines: %w", err)
	}
	return w.History, nil
}

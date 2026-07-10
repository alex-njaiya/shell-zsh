package utils

import (
	"os"
	"strings"
)

func CompletionsFromFiles(prefix string) []string {
	var matches []string

	currentDir, err := os.Getwd()

	if err != nil {
		return matches
	}

	// loop through all paths
	entries, err := os.ReadDir(currentDir)

	if err != nil {
		return matches
	}

	for _, e := range entries {
		if strings.HasPrefix(e.Name(), prefix) {
			matches = append(matches, e.Name())
		}
	}

	return matches
}

func LongestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]

	// iterate over the str string by strings
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			prefix = prefix[:len(prefix)-1]
		}
	}

	return prefix
}

func CompletionFromPath(prefix string) []string {
	var matches []string

	paths := os.Getenv("PATH")

	// split the paths using the semicolon
	dirs := strings.Split(paths, ":")

	// iterate over the dirs and get into each directory
	for _, dir := range dirs {
		// for every directory check if the entry matches the prefix
		// loop over each directory and match its name with the prefix
		dirContents, err := os.ReadDir(dir)

		if err != nil {
			continue
		}

		// check each dirEntry content whether they match with the prefix
		for _, entry := range dirContents {
			if strings.HasPrefix(entry.Name(), prefix) {
				//apend to the matches slice
				matches = append(matches,  entry.Name())
			}
		}
	}

	return matches
}

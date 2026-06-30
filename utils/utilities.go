package utils

import (
	"os"
	"strings"
)

func CompletionsFromPath(prefix string) []string {
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

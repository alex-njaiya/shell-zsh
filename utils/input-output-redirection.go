package utils

import (
	"fmt"
)

func InputOutputRedirection(token []string) {
	// read the token and check for the symbols: >, >> for writing to a new file that does not exist
	// read the token and check for the symbols: <, << 

	for i, symbol := range token {
		if symbol == ">" {
			fmt.Println("This is input redirection")
			// get the part that is the command and the part that is the argument
			command := token[:i - 1]
			argument := token[i + 1:]

			// execute the command and also the argument
			// for the argument open the file for read and write operation only and write to that file

			fmt.Println(command)
			fmt.Println(argument)
		}
	}
}
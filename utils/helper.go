package utils

import (
	"fmt"
	"os"
	"strings"
)

func ExpandVariables(args []string) []string {
	for i, token := range args {
		if strings.HasPrefix(token, "$") {
			if len(token) > 1 {

				val := token[1:]
				delimeter := "/"
				// if the token is just a pure variable reference
				parts := strings.SplitN(val, delimeter, 2)
				varName := parts[0]

				if len(parts) == 2 {
					suffix := parts[1]

					// The token is a variable with path prefix
					// Reconstruct the token with the "/"
					env, found := os.LookupEnv(varName)

					if !found { // Not found
						args[i] = "/" + suffix
					}

					// If the environment variable has been found

					// reconstruct the env with the entire path
					args[i] = env + "/" + suffix
				}

				if len(parts) == 1 { // Meaning it is a pure variable with no path reference
					env, found := os.LookupEnv(varName)

					if !found {
						args[i] = ""
					}

					args[i] = env
				}

			}

		}

		if token == "$$" {
			// get the current process id
			currentPID := os.Getegid()
			fmt.Println(currentPID)
		}
		if strings.HasPrefix(token, "${") && strings.HasSuffix(token, "}") {
			// extract the variable name
			// I am thinking maybe I should remove the first 2 characters at the beginning and one at the end

			parts := strings.SplitN(token, "/", 2)
			varName := parts[0]

			variableName := varName[2 : len(varName)-1]

			if len(parts) == 2 {
				// If the length of parts is 2 we need to reconstruct the path
				suffix := parts[1]

				// lookup if the env is found
				env, found := os.LookupEnv(variableName)

				if !found {
					args[i] = "/" + suffix
				}

				args[i] = env + "/" + suffix
			} 

			env, found := os.LookupEnv(variableName)

			if !found {
				args[i] = ""
			}

			args[i] = env


			// if there is something after the }(the last closing bracket) reconstruct the path
		}
	}

	return args
}

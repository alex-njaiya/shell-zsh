package utils

import (
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
	}

	return args
}

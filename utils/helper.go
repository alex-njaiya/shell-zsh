package utils

import (
	"os"
	"strings"
)

func ExpandVariables(args []string) []string {
	for i, token := range args {
		if strings.HasPrefix(token, "$") {
			val := token[1:]

			if len(token) > 1 {
				// strip everything that is not a letter, digit or underscore
				extract := func(r rune) bool {
					return r == '/'
				}

				parts := strings.FieldsFunc(val, extract)

				env, ok := os.LookupEnv(parts[0])

				if !ok {
					args[i] = " "
				}

				// recontruct the token if it had a suffix
				if len(parts) == 1 {
					// no need to reconstruct
					args[i] = env
				} else {
					// there are more than 2 parts lets add the suffix
					result := strings.Join(parts[1:], "/")
					args[i] = result
				}

			}
		}
	}

	return args
}

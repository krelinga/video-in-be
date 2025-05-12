package env

import (
	"fmt"
	"os"
)

func getString(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("Environment variable '%s' not set", key))
	}
	return value
}
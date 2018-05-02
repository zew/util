package util

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Prime data source is in the "dsn1" configuration map key.
// Except it is set by ENV var "DATASOURCE1".
func PrimeDataSource() string {
	dsn1 := os.Getenv("DATASOURCE1")
	if dsn1 == "" {
		dsn1 = "dsn1"
	}
	return dsn1
}

// Required environment var terminates the program
// if the env var is not set.
func EnvVarRequired(key string) string {
	all := os.Environ() // slice with key=val entries
	found := false
	compareTo := key + "="
	for _, v := range all {
		if strings.HasPrefix(v, compareTo) {
			found = true
			break
		}
	}
	if !found {
		log.Printf("----")
		log.Printf("Program *requires* environment variable %q.\nExiting.", key)
		log.Printf("----")
		os.Exit(1)
	}

	return os.Getenv(key)
}

// EnvVar returns an error, if the key is not present
func EnvVar(key string) (string, error) {
	all := os.Environ() // slice with key=val entries
	found := false
	compareTo := key + "="
	for _, v := range all {
		if strings.HasPrefix(v, compareTo) {
			found = true
			break
		}
	}
	if found {
		return os.Getenv(key), nil
	}
	return "", fmt.Errorf("ENV variable %v does not exist", key)
}

package util

import (
	"log"
	"os"
	"strings"
)

func PrimeDataSource() string {
	dsn1 := os.Getenv("DATASOURCE1")
	if dsn1 == "" {
		dsn1 = "dsn1"
	}
	return dsn1
}

func EnvVar(key string) string {
	all := os.Environ()
	found := false
	for _, v := range all {
		if strings.HasPrefix(v, key) {
			found = true
		}
	}
	if !found {
		log.Printf("----")
		log.Printf("Program *requires* environment variable %q.\nExiting.", key)
		log.Printf("----")
		os.Exit(1)
	}

	envVal := os.Getenv(key)
	return envVal
}

package env

import (
	"log"
	"os"
	"strconv"
)

func GetString(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Missing environment variable:", key)
	}

	return val
}

func GetInt(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalln("Missing environment variable:", key)
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalln("Cannot parse environment variable:", key)
	}

	return valAsInt
}

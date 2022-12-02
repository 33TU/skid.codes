package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// TODO: find better solution.
var (
	testing       = strings.HasSuffix(os.Args[0], ".test")
	testingEnv    = "../.env"
	productionEnv = ".env"
)

// init loads enviroment variables from .env file.
func init() {
	err := godotenv.Load(getEnvPath())

	if err != nil {
		log.Fatalln(err)
	}
}

// getEnvPath get .env file's path
func getEnvPath() string {
	if testing {
		return testingEnv
	}

	return productionEnv
}

// Get finds and returns an environment variable and a boolean indicating whether the variable was found.
func Get(key string) (string, bool) {
	if env, ok := os.LookupEnv(key); ok {
		return env, true
	}

	return "", false
}

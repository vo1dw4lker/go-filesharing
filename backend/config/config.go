package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Environment struct {
	PostgresUser      string        // POSTGRES_USER
	PostgresPassword  string        // POSTGRES_PASSWORD
	PostgresDB        string        // POSTGRES_DB
	AllowedExpOptions []int         // APP_ALLOWED_EXP_OPTIONS
	DbTimeout         time.Duration // APP_DB_TIMEOUT
	StorageDir        string        // APP_STORAGE_DIR
}

var Env = ParseEnv()

func ParseEnv() *Environment {
	env := &Environment{}
	for _, key := range os.Environ() {
		keyVal := strings.Split(key, "=")
		if len(keyVal) == 2 {
			envField, envValue := keyVal[0], keyVal[1]
			envValue = unquote(envValue)

			switch envField {
			case "POSTGRES_USER":
				env.PostgresUser = envValue
			case "POSTGRES_PASSWORD":
				env.PostgresPassword = envValue
			case "POSTGRES_DB":
				env.PostgresDB = envValue
			case "APP_ALLOWED_EXP_OPTIONS":
				err := env.parseExpOptions(envValue)
				if err != nil {
					log.Fatalln("Failed parsing APP_ALLOWED_EXP_OPTIONS:", err)
				}
			case "APP_DB_TIMEOUT":
				err := env.parseDBTimeout(envValue)
				if err != nil {
					log.Fatalln("Failed parsing APP_DB_TIMEOUT: ", err)
				}
			case "APP_STORAGE_DIR":
				env.StorageDir = envValue
			}
		}
	}
	return env
}

func (env *Environment) parseDBTimeout(envValue string) error {
	timeout, err := strconv.Atoi(envValue)
	if err != nil {
		return err
	}

	env.DbTimeout = time.Second * time.Duration(timeout)
	return err
}

func (env *Environment) parseExpOptions(envValue string) error {
	var options []string
	// Parse the string to a slice of ints
	options = strings.Split(envValue, ", ")

	converted := make([]int, len(options))
	// Convert each string to an int
	for i, option := range options {
		if value, err := strconv.Atoi(option); err != nil {
			return err
		} else {
			converted[i] = value
		}
	}

	env.AllowedExpOptions = converted
	return nil
}

func unquote(s string) string {
	return strings.ReplaceAll(s, "'", "")
}

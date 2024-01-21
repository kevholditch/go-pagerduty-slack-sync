package sync

import (
	"os"
	"strconv"
)

func GetEnvStr(key string, d string) string {
	v := os.Getenv(key)
	if v == "" {
		return d
	}
	return v
}

func GetEnvInt(key string, d int) int {
	s, ok := os.LookupEnv(key)
	if !ok {
		return d
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return v
}

func GetEnvBool(key string, d bool) bool {
	s, ok := os.LookupEnv(key)
	if !ok {
		return d
	}

	v, err := strconv.ParseBool(s)
	if err != nil {
		return d
	}
	return v
}

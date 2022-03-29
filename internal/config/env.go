package config

import "os"

func GetXORKey() string {
	result := os.Getenv("XOR_KEY")
	if result == "" {
		result = "simplekey"
	}
	return result
}

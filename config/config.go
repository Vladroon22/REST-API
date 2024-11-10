package config

import "os"

type Config struct {
	Addr_PORT string
	DB        string
	JwtKEY    string
}

func CreateConfig() *Config {
	return &Config{
		Addr_PORT: getEnv("addr_port", ":3000"),
		DB:        getEnv("DB", ""),
		JwtKEY:    getEnv("KEY", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

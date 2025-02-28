package configs

import (
	"os"
)

type Config struct {
	User     string
	Password string
	Dbname   string
	Sslmode  string
	Host     string
	Port     string
	ApiKey   string
}

func New() *Config {
	return &Config{
		User:     getEnv("user", ""),
		Password: getEnv("password", ""),
		Dbname:   getEnv("dbname", ""),
		Sslmode:  getEnv("sslmode", ""),
		Host:     getEnv("host", ""),
		Port:     getEnv("port", ""),
		ApiKey:   getEnv("apikey", ""),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

//func getEnvAsInt(name string, defaultVal int) int {
//	valueStr := getEnv(name, "")
//	if value, err := strconv.Atoi(valueStr); err == nil {
//		return value
//	}
//
//	return defaultVal
//}

package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	DatabaseURL   string
	MQTTBrokerURL string
	MQTTTopic     string
	HTTPPort      string
}

func LoadConfig() *Config {
	loadDotEnv("../.env")
	loadDotEnv(".env")

	cfg := &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://agriedge:agriedge@localhost:5432/agriedge?sslmode=disable"),
		MQTTBrokerURL: getEnv("MQTT_BROKER_URL", "tcp://localhost:1883"),
		MQTTTopic:     getEnv("MQTT_TOPIC", "agriedge/+/daily_metrics"),
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// loadDotEnv simple helper to parse .env file without external library
func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			if os.Getenv(k) == "" {
				os.Setenv(k, v)
			}
		}
	}
}

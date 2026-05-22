package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddress  string
	ServerPort     string
	Mode           string

	PVEEndpoint    string
	PVEUsername    string
	PVEPassword    string
	PVERealm       string
	PVETLSSkip     bool
	PVETimeout     int

	MaxConcurrentClones int
	PVEWorkers          int

	NodePollInterval    time.Duration
	ErrorCleanupMaxAge  time.Duration
	ErrorAlertThreshold float64

	DefaultClassroomID    string
	DefaultTotalStudents  int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress:         getEnv("SERVER_ADDRESS", "0.0.0.0"),
		ServerPort:            getEnv("SERVER_PORT", "8080"),
		Mode:                  getEnv("GIN_MODE", "release"),

		PVEEndpoint:           getEnv("PVE_ENDPOINT", "https://pve-cluster.local:8006"),
		PVEUsername:           getEnv("PVE_USERNAME", "root@pam"),
		PVEPassword:           getEnv("PVE_PASSWORD", ""),
		PVERealm:              getEnv("PVE_REALM", "pam"),
		PVETLSSkip:            getEnvBool("PVE_TLS_SKIP_VERIFY", true),
		PVETimeout:            getEnvInt("PVE_TIMEOUT_SECONDS", 30),

		MaxConcurrentClones:   getEnvInt("MAX_CONCURRENT_CLONES", 50),
		PVEWorkers:            getEnvInt("PVE_WORKERS", 10),

		NodePollInterval:      getEnvDuration("NODE_POLL_INTERVAL", 5),
		ErrorCleanupMaxAge:    getEnvDuration("ERROR_CLEANUP_MAX_AGE", 30),
		ErrorAlertThreshold:   getEnvFloat64("ERROR_ALERT_THRESHOLD", 20.0),

		DefaultClassroomID:    getEnv("DEFAULT_CLASSROOM_ID", "classroom-01"),
		DefaultTotalStudents:  getEnvInt("DEFAULT_TOTAL_STUDENTS", 52),

		ReadTimeout:           60 * time.Second,
		WriteTimeout:          60 * time.Second,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func getEnvFloat64(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return f
}

func getEnvDuration(key string, defaultValueMinutes int) time.Duration {
	minutes := getEnvInt(key, defaultValueMinutes)
	return time.Duration(minutes) * time.Minute
}

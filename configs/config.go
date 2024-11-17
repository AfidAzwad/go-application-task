package configs

import (
	"os"
)

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// Read and Write config to distribute load

func GetWriteDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("WRITE_DB_HOST"),
		Port:     5432,
		User:     os.Getenv("WRITE_DB_USER"),
		Password: os.Getenv("WRITE_DB_PASSWORD"),
		DBName:   os.Getenv("WRITE_DB_NAME"),
	}
}

func GetReadDBConfig() DBConfig {
	return DBConfig{
		Host:     os.Getenv("READ_DB_HOST"),
		Port:     5432,
		User:     os.Getenv("READ_DB_USER"),
		Password: os.Getenv("READ_DB_PASSWORD"),
		DBName:   os.Getenv("READ_DB_NAME"),
	}
}

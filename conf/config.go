package conf

import (
	"flag"
	"os"
)

type Config struct {
	RunAddress     string
	AccrualAddress string
	SecretKey      string
	Database       *DatabaseConfig
}

type DatabaseConfig struct {
	Driver                  string
	DatabaseURI             string
	ConnMaxLifeTimeInMinute int
	MaxOpenConns            int
	MaxIdleConns            int
}

/*
	адрес и порт запуска сервиса: переменная окружения ОС RUN_ADDRESS или флаг -a;
	адрес подключения к базе данных: переменная окружения ОС DATABASE_URI или флаг -d;
	адрес системы расчёта начислений: переменная окружения ОС ACCRUAL_SYSTEM_ADDRESS или флаг -r.
*/

func NewConfig() *Config {
	runAddress := flag.String("a", "", "Адрес и порт запуска сервиса")
	databaseURI := flag.String("d", "", "Адрес подключения к базе данных")
	accrualAddress := flag.String("r", "", "Адрес системы расчёта начислений")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		*runAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		*databaseURI = envDatabaseURI
	}

	if envAccrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddress != "" {
		*accrualAddress = envAccrualAddress
	}

	var secretKey string
	if envSecretKey := os.Getenv("SECRET_KEY"); envSecretKey != "" {
		secretKey = envSecretKey
	}

	return &Config{
		RunAddress:     *runAddress,
		AccrualAddress: *accrualAddress,
		SecretKey:      secretKey,
		Database: &DatabaseConfig{
			Driver:                  "pgx",
			DatabaseURI:             *databaseURI,
			ConnMaxLifeTimeInMinute: 3,
			MaxOpenConns:            10,
			MaxIdleConns:            1,
		},
	}
}

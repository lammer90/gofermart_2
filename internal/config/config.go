package config

import (
	"flag"
	"os"
)

var ServAddress string
var DataSource string
var PrivateKey string
var AccrualAddress string

func InitConfig() {
	initFlags()
	initEnv()
}

func initFlags() {
	flag.StringVar(&ServAddress, "a", ":8080", "Request URL")
	flag.StringVar(&DataSource, "d", "postgresql://localhost:5432/plotnikov?user=postgres&password=1234", "DataSource path")
	flag.StringVar(&PrivateKey, "p", "privateKey", "PrivateKey for jwt auth")
	flag.StringVar(&AccrualAddress, "r", "http://localhost:8081", "Accrual address")
	flag.Parse()
}

func initEnv() {
	if envServAddr := os.Getenv("RUN_ADDRESS"); envServAddr != "" {
		ServAddress = envServAddr
	}

	if envDataSource := os.Getenv("DATABASE_URI"); envDataSource != "" {
		DataSource = envDataSource
	}

	if privateKey := os.Getenv("PRIVATE_KEY"); privateKey != "" {
		PrivateKey = privateKey
	}

	if accrualAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); accrualAddress != "" {
		AccrualAddress = accrualAddress
	}
}

package main

import (
	"fmt"
	"time"

	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/currency"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ No .env file found, using system environment variables")
	}

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	currency.CheckEuro()

	for range ticker.C {
		currency.CheckEuro()
	}
}

package main

import (
	"fmt"
	"os"
	"time"

	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/currency"
	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/email"
	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/route"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️ No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DB_CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(fmt.Errorf("could not open database: %w", err))
	}

	db.AutoMigrate(&currency.CurrencyAlarm{})

	emailSender := email.NewEmailSender()
	currencyChecker := currency.NewCurrencyChecker()
	checkCurrencyUC := currency.NewCheckCurrencyUC(emailSender, currencyChecker, db)
	createCurrencyAlarmUC := currency.NewCreateCurrencyAlarmUC(db)
	deleteCurrencyAlarmUC := currency.NewDeleteCurrencyAlarmUC(db)

	go StartAlarmWorker(db, *checkCurrencyUC)

	e := echo.New()
	e.HTTPErrorHandler = route.ProblemDetailHTTPErrorHandler

	routers := []route.Router{
		currency.NewCurrencyAlarmRouter(*createCurrencyAlarmUC, *deleteCurrencyAlarmUC),
	}

	for _, router := range routers {
		router.Route(e)
	}

	e.Logger.Fatal(e.Start(":8080"))
}

func StartAlarmWorker(database *gorm.DB, checkCurrencyUC currency.CheckCurrencyUC) {
	for {
		var currencyAlarms []currency.CurrencyAlarm
		result := database.Find(&currencyAlarms)
		if result.Error != nil {
			fmt.Printf("🔴 Error trying to fetch the alarms")
			return
		}

		for _, currencyAlarm := range currencyAlarms {
			go checkCurrencyUC.Execute(&currencyAlarm)
		}

		time.Sleep(1 * time.Hour)
	}
}

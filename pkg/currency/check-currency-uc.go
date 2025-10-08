package currency

import (
	"bytes"
	"database/sql"
	"fmt"
	"html/template"
	"time"

	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/email"
	"github.com/EduardoPPCaldas/CurrencyAlarm/templates"
	"gorm.io/gorm"
)

type CheckCurrencyUC struct {
	emailSender     email.EmailSender
	currencyChecker CurrencyChecker
	db              *gorm.DB
}

func NewCheckCurrencyUC(emailSender email.EmailSender, currencyChecker CurrencyChecker, db *gorm.DB) *CheckCurrencyUC {
	return &CheckCurrencyUC{
		emailSender:     emailSender,
		currencyChecker: currencyChecker,
		db:              db,
	}
}

func (uc CheckCurrencyUC) Execute(currencyAlarm *CurrencyAlarm) error {
	bid, err := uc.currencyChecker.CheckCurrency(currencyAlarm.OwnedCurrency, currencyAlarm.ConvertedCurrency)
	if err != nil {
		return err
	}

	if bid > currencyAlarm.Threshold {
		fmt.Printf("Currency above Threshold")
		return nil
	}

	if currencyAlarm.AlarmedAt.Valid && currencyAlarm.AlarmedAt.Time.Truncate(24 * time.Hour).Equal(time.Now().UTC().Truncate(24 * time.Hour)) {
		fmt.Printf("Currency Already alarmed")
		return nil
	}

	fmt.Printf("🚨 Currency dropped! Current EUR/BRL: %.2f\n", bid)

	subject := fmt.Sprintf("Currency Alarm: %s < %.2f", currencyAlarm.ConvertedCurrency, currencyAlarm.Threshold)

	tmpl, err := template.ParseFS(templates.FS, "currency-alarm.html")
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, map[string]any{"CurrencyAlarm": currencyAlarm, "Bid": bid}); err != nil {
		return err
	}

	msg := []byte(fmt.Sprintf(
		"Subject: %s\r\n"+
			"MIME-version: 1.0;\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n"+
			"%s",
		subject, &body,
	))

	if err := uc.emailSender.SendEmail([]string{currencyAlarm.Email}, msg); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("✅ Alert email sent!")

	currencyAlarm.AlarmedAt = sql.NullTime{
		Time:  time.Now().UTC(),
		Valid: true,
	}

	uc.db.Save(&currencyAlarm)

	return nil
}

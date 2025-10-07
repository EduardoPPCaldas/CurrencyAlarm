package currency

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/EduardoPPCaldas/CurrencyAlarm/pkg/email"
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

	if currencyAlarm.AlarmedAt.Valid && currencyAlarm.AlarmedAt.Time.Format(time.DateOnly) == time.Now().UTC().Format(time.DateOnly) {
		fmt.Printf("Currency Already alerted")
		return nil
	}

	fmt.Printf("🚨 Currency dropped! Current EUR/BRL: %.2f\n", bid)

	subject := fmt.Sprintf("Currency Alert: %s < %.2f", currencyAlarm.ConvertedCurrency, currencyAlarm.Threshold)
	body := fmt.Sprintf("The Currency just dropped below %.2f BRL!\n\nCurrent rate: %.2f", currencyAlarm.Threshold, bid)

	if err := uc.emailSender.SendEmail(subject, []string{currencyAlarm.Email}, []byte(body)); err != nil {
		fmt.Println("Error sending email:", err)
		return err
	}

	fmt.Println("✅ Alert email sent!")

	currencyAlarm.AlarmedAt = sql.NullTime{
		Time: time.Now().UTC(),
		Valid: true,
	}

	uc.db.Save(&currencyAlarm)

	return nil
}

package currency

import (
	"database/sql"

	"github.com/google/uuid"
)

type CurrencyAlarm struct {
	ID                uuid.UUID
	Email             string
	OwnedCurrency     string
	ConvertedCurrency string
	Threshold         float64
	AlarmedAt         sql.NullTime
}

func NewCurrencyAlarm(email, ownedCurrency, convertedCurrency string, threshold float64) *CurrencyAlarm {
	return &CurrencyAlarm{
		ID:                uuid.New(),
		Email:             email,
		OwnedCurrency:     ownedCurrency,
		ConvertedCurrency: convertedCurrency,
		Threshold:         threshold,
		AlarmedAt: sql.NullTime{
			Valid: false,
		},
	}
}

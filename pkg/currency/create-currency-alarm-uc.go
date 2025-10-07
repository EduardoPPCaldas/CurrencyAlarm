package currency

import (
	"context"

	"gorm.io/gorm"
)

type CreateCurrencyAlarmDto struct {
	Email             string
	OwnedCurrency     string
	ConvertedCurrency string
	Threshold         float64
}

type CreateCurrencyAlarmUC struct {
	db *gorm.DB
}

func NewCreateCurrencyAlarmUC(db *gorm.DB) *CreateCurrencyAlarmUC {
	return &CreateCurrencyAlarmUC{
		db: db,
	}
}

func (uc CreateCurrencyAlarmUC) Execute(ctx context.Context, dto CreateCurrencyAlarmDto) error {
	currencyAlarm := NewCurrencyAlarm(dto.Email, dto.OwnedCurrency, dto.ConvertedCurrency, dto.Threshold)

	err := gorm.G[CurrencyAlarm](uc.db).Create(ctx, currencyAlarm)
	if err != nil {
		return err
	}

	return nil
}

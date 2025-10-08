package currency

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeleteCurrencyAlarmUC struct {
	db *gorm.DB
}

func NewDeleteCurrencyAlarmUC(db *gorm.DB) *DeleteCurrencyAlarmUC {
	return &DeleteCurrencyAlarmUC{
		db: db,
	}
}

func (uc DeleteCurrencyAlarmUC) Execute(ctx context.Context, id uuid.UUID) error {
	row, err := gorm.G[CurrencyAlarm](uc.db).Where("id = ?", id).Delete(ctx)
	if err != nil {
		return err
	}

	if row == 0 {
		return fmt.Errorf("no row were affected")	
	}

	return nil
}
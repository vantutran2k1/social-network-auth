package transaction

import (
	"gorm.io/gorm"
)

type TransactionManager struct {
	DB *gorm.DB
}

var TxManager *TransactionManager

func InitTransactionManager(db *gorm.DB) {
	TxManager = &TransactionManager{
		DB: db,
	}
}

func (t *TransactionManager) WithTransaction(fn func(tx *gorm.DB) error) error {
	tx := t.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

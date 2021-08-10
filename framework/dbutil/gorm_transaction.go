package dbutil

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// TransactionTask 事务处理函数
type TransactionTask func(db *gorm.DB) error

// Transaction 事务
func Transaction(db *gorm.DB, tasks ...TransactionTask) (rerr error) {
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return errors.Wrap(err, "db transaction start failed")
	}

	defer func() {
		if e := recover(); e != nil {
			tx.Rollback()
			rerr = errors.Errorf("db transaction error: %v", e)
		}
	}()

	for _, task := range tasks {
		if err := task(tx); err != nil {
			if rerr := tx.Rollback().Error; rerr != nil {
				return errors.Wrap(rerr, "db transaction rollback failed")
			}
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "db transaction commit failed")
	}
	return nil
}

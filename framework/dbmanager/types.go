package dbmanager

import "gorm.io/gorm"

type IDatabaseManager interface {
	Register(name string, dbc *gorm.DB) error
	Get(name string) (*gorm.DB, error)
}

package dberror

import (
	"github.com/go-sql-driver/mysql"
)

func IsDuplicateEntry(err error) bool {
	//if sqliteErr, ok := err.(sqlite3.Error); ok {
	//	return sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique ||
	//		sqliteErr.ExtendedCode == sqlite3.ErrConstraintPrimaryKey
	//} else
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return mysqlErr.Number == errorCodeMySQLDuplicateEntry ||
			mysqlErr.Number == errorCodeMySQLDuplicateEntryWithKeyName
	}

	return false
}

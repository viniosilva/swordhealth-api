package repository

type MySQLErrorCode int

const (
	MySQLErrorCodeForeignKeyConstraint MySQLErrorCode = 1452
)

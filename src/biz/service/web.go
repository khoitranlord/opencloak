package service

import (
	"github.com/Masterminds/squirrel"
	custdb "github.com/khoitranlord/opencloak/src/internal/db"
)

type WebService struct {
	db         *custdb.LayeredDb
	sqlBuilder squirrel.StatementBuilderType
}

func NewWebService() *WebService {
	return &WebService{
		db:         custdb.Layered(),
		sqlBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

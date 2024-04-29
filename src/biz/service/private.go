package service

import (
	custdb "github.com/khoitranlord/opencloak/src/internal/db"
)

type PrivateService struct {
	db         *custdb.LayeredDb
	webService *WebService
}

func NewPrivateService(webService *WebService) *PrivateService {
	return &PrivateService{
		db:         custdb.Layered(),
		webService: webService,
	}
}

package service

import (
	"context"
	"sync"

	"github.com/khoitranlord/opencloak/src/internal/configs"
)

var once sync.Once

var (
	privateService *PrivateService
	webService     *WebService
)

func Init(c *configs.Configs, ctx context.Context) {
	once.Do(func() {
		webService = NewWebService()
		privateService = NewPrivateService(
			webService)
	})
}

func GetWebService() *WebService {
	return webService
}

func GetPrivateService() *PrivateService {
	return privateService
}

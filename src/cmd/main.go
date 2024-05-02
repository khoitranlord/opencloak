package main

import (
	"context"
	"time"

	publicapi "github.com/khoitranlord/opencloak/src/api/public"
	"github.com/khoitranlord/opencloak/src/biz/service"
	"github.com/khoitranlord/opencloak/src/internal/app"
	"github.com/khoitranlord/opencloak/src/internal/configs"
	custdb "github.com/khoitranlord/opencloak/src/internal/db"
	custhttp "github.com/khoitranlord/opencloak/src/internal/http"
	"github.com/khoitranlord/opencloak/src/internal/logger"
	"github.com/khoitranlord/opencloak/src/models/db"

	"go.uber.org/zap"
)

func main() {
	globalCtx, cancelGlobalCtx := context.WithCancel(context.Background())

	app.Run(
		time.Second*10,
		func(configs *configs.Configs, zl *zap.Logger) []app.Optioner {
			return []app.Optioner{
				app.WithHttpServer(custhttp.New(
					custhttp.WithGlobalConfigs(&configs.Public),
					custhttp.WithRegistration(publicapi.ServiceRegistration()),
					custhttp.WithMiddleware(custhttp.CommonPublicMiddlewares(&configs.Public)...), // Fixed parameter passing
				)),
				app.WithFactoryHook(func() error {
					custdb.Init(globalCtx, configs)
					err := custdb.Migrate(custdb.Gorm(), &db.User{}) // Assuming db.Users is defined somewhere
					if err != nil {
						return err
					}
					service.Init(configs, globalCtx)
					return nil // Added return statement
				}),
				app.WithShutdownHook(func(ctx context.Context) {
					cancelGlobalCtx()
					custdb.Stop(ctx)
					logger.Close()
				}),
			}
		},
	)
}

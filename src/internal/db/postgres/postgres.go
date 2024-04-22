package custpg

import (
	"context"

	"github.com/jmoiron/sqlx"
	config "github.com/khoitranlord/opencloak/src/internal/configs"
	custerror "github.com/khoitranlord/opencloak/src/internal/error"
	"github.com/khoitranlord/opencloak/src/internal/logger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Options struct {
	globalConfigs *config.DatabaseConfigs
}

type Optioner func(*Options)

func WithConfigs(globalConfigs *config.DatabaseConfigs) Optioner {
	return func(o *Options) {
		o.globalConfigs = globalConfigs
	}
}

func NewSqlx(ctx context.Context, options ...Optioner) (*sqlx.DB, error) {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	globalConfigs := opts.globalConfigs

	logger.SInfo("db.postgres.Init: create SQLx connection",
		zap.String("connection", globalConfigs.Connection))

	client, err := sqlx.Connect("postgres", globalConfigs.Connection)
	if err != nil {
		logger.SFatal("db.postgres.Init: open error",
			zap.Error(err))
		return nil, err
	}

	return client, nil
}

func NewGorm(ctx context.Context, options ...Optioner) (*gorm.DB, error) {
	opts := &Options{}
	for _, opt := range options {
		opt(opts)
	}

	connString := opts.globalConfigs.Connection
	db, err := gorm.Open(
		postgres.Open(connString),
		&gorm.Config{})
	if err != nil {
		return nil, custerror.FormatInternalError("buildGorm: err = %s", err)
	}

	db.Exec("CREATE EXTENSION IF NOT EXISTS vector")

	return db, nil
}

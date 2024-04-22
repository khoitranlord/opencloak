package custhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	custerror "github.com/khoitranlord/opencloak/src/internal/error"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3"
	"github.com/khoitranlord/opencloak/src/internal/configs"
)

type RegistrationFunc func(app *gin.Engine)

type HttpServer struct {
	configs *HttpServerConfigs
	app     *http.Server
}

type HttpServerConfigs struct {
	configs      *configs.HttpConfigs
	registration RegistrationFunc
	errorHandler fiber.ErrorHandler
	middlewares  []gin.HandlerFunc
	templatePath string
}

type Optioner func(config *HttpServerConfigs)

func WithErrorHandler(handler fiber.ErrorHandler) Optioner {
	return func(config *HttpServerConfigs) {
		config.errorHandler = handler
	}
}

func WithGlobalConfigs(conf *configs.HttpConfigs) Optioner {
	return func(configs *HttpServerConfigs) {
		configs.configs = conf
	}
}

func WithRegistation(handler RegistrationFunc) Optioner {
	return func(configs *HttpServerConfigs) {
		configs.registration = handler
	}
}

func WithMiddleware(middlewares ...gin.HandlerFunc) Optioner {
	return func(configs *HttpServerConfigs) {
		configs.middlewares = middlewares
	}
}

func WithTemplatePath(path string) Optioner {
	return func(configs *HttpServerConfigs) {
		configs.templatePath = path
	}
}

func (s *HttpServer) Start() error {
	globalConfigs := s.configs.configs
	tlsConfigs := globalConfigs.Tls
	port := fmt.Sprintf(":%d", globalConfigs.Port)
	s.app.Addr = port

	if tlsConfigs.Enabled() {
		if err := s.app.ListenAndServeTLS(tlsConfigs.Cert, tlsConfigs.Key); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return custerror.FormatInternalError("HttpServer.Start: err = %s", err)
		}
	} else {
		if err := s.app.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return err
			}
			return custerror.FormatInternalError("HttpServer.Start: err = %s", err)
		}
	}

	return nil

}

func (s *HttpServer) Name() string {
	return s.configs.configs.Name
}

func (s *HttpServer) Stop(ctx context.Context) error {
	globalConfigs := s.configs.configs

	if err := s.app.Shutdown(ctx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return custerror.FormatTimeout("HttpServer.Stop: server stopping deadline exceeded name = %s", globalConfigs.Name)
		}
		return custerror.FormatInternalError("HttpServer.Shutdown: err = %s", err)
	}

	return nil
}

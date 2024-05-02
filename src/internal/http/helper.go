package custhttp

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/khoitranlord/opencloak/src/internal/configs"
	custerror "github.com/khoitranlord/opencloak/src/internal/error"
	"github.com/khoitranlord/opencloak/src/internal/logger"
	"net/http"
	"time"
)

func CommonPublicMiddlewares(configs *configs.HttpConfigs) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		CORSMiddleware(),
		ginzap.Ginzap(logger.Logger(), time.RFC3339, true),
		gzip.Gzip(gzip.BestCompression),
		gin.Recovery(),
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func ToHTTPErr(err error, c *gin.Context) {
	customError, yes := err.(*custerror.CustomError)
	if yes {
		customError.Gin(c)
		return
	}
	madeCustom := custerror.NewError(
		err.Error(),
		http.StatusInternalServerError)
	madeCustom.Gin(c)
}

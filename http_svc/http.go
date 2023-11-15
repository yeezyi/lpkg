package http_svc

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/yeezyi/gin-swagger"
	"net/http"
)

type Config struct {
	Addr        string
	Router      func(router gin.IRouter)
	SwaggerOpen bool
}

func NewServer(config *Config) *http.Server {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(cors())

	swaggerRouter(config.SwaggerOpen, engine)
	config.Router(engine)
	var addr = ":8080"
	if len(config.Addr) != 0 {
		addr = config.Addr
	}
	hs := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	return hs
}

func cors() gin.HandlerFunc {
	const (
		AccessControlAllowOrigin      = "Access-Control-Allow-Origin"
		AccessControlMaxAge           = "Access-Control-Max-Age"
		AccessControlAllowMethods     = "Access-Control-Allow-Methods"
		AccessControlAllowHeaders     = "Access-Control-Allow-Headers"
		AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	)
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set(AccessControlAllowOrigin, "*")
		ctx.Writer.Header().Set(AccessControlMaxAge, "86400")
		ctx.Writer.Header().Set(AccessControlAllowMethods, "*")
		ctx.Writer.Header().Set(AccessControlAllowHeaders, "*")
		ctx.Writer.Header().Set(AccessControlAllowCredentials, "true")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
		ctx.Next()
	}
}

func swaggerRouter(open bool, r gin.IRouter) {
	router := r.Group("/swagger")
	{
		if open {
			router.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.DocExpansion("none"), ginSwagger.DefaultModelsExpandDepth(10)))
		} else {
			//router.GET("/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "ems.swagger_web_close"))
			router.GET("/*any", func(c *gin.Context) {
				c.String(http.StatusNotFound, "")
			})
		}
	}
}

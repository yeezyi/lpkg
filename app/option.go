package app

import (
	"github.com/gin-gonic/gin"
	"github.com/yeezyi/lpkg/config"
	"github.com/yeezyi/lpkg/log"
	"google.golang.org/grpc"
)

/*
var application *app.App
application = app.Init(&app.Config{
Name:   "无敌~",
Router: router.Router,
})

if err := application.Execute(); err != nil {
log.Error(err)
return
}
*/

type ConfigItem struct {
	Source string
	Cfg    config.Cfg
}

type Config struct {
	Name string

	Logger log.Logger

	ConfigSource []ConfigItem

	Router func(router gin.IRouter)

	GrpcRegister func(*grpc.Server)

	Init func()
}

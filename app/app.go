package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	redisLib "github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	logrusLib "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/yeezyi/lpkg/config"
	"github.com/yeezyi/lpkg/config/viper"
	"github.com/yeezyi/lpkg/http_svc"
	"github.com/yeezyi/lpkg/log"
	"github.com/yeezyi/lpkg/log/logrus"
	"github.com/yeezyi/lpkg/repository/gorm"
	"github.com/yeezyi/lpkg/repository/redis"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	gormLib "gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
)

type App struct {
	cmd      *cobra.Command
	DB       *gormLib.DB
	Redis    *redisLib.Client
	http     *http.Server
	grpc     *grpc.Server
	grpcAddr string

	//cfg *Config
}

//var Std *App

const (
	ConfigSourceFile = "std"

	MysqlDsn         = "dsn"
	MysqlMaxIdleConn = "max-idle-conns"
	MysqlMaxOpenConn = "max-open-conns"
	MysqlMaxLifeTime = "max-lifetime"

	RedisAddr     = "addr"
	RedisPassword = "password"
	RedisDb       = "db"
	RedisPoolsize = "poolsize"

	ArgLogLevel  = "log.level"
	ArgLogCaller = "log.caller"
	ArgLogFormat = "log.format"
	ArgLogPretty = "log.pretty"

	SrvHttpAddr    = "srv.http-addr"
	SrvGrpcAddr    = "srv.grpc-addr"
	SrvSwaggerOpen = "srv.swagger-open"
)

func init() {
	log.SetLogger(logrus.NewLogger())
}

var app *App

func GetInstance() *App {
	return app
}

func Init(cfg *Config) *App {
	app = new(App)
	cmd := &cobra.Command{
		Use:           cfg.Name,
		SilenceUsage:  true,
		SilenceErrors: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return app.postRun(cfg)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return app.run()
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			app.stop()
		},
	}
	app.cmd = cmd

	cmd.Flags().String(ArgLogLevel, "info", "可选日志级别：debug/info/warn/error")
	cmd.Flags().Bool(ArgLogCaller, false, "是否显示调用方法：true/false")
	cmd.Flags().String(ArgLogFormat, "text", "日志格式：text/json")
	cmd.Flags().Bool(ArgLogPretty, false, "json格式日志打印美化：true/false")

	return app
}
func (app *App) postRun(cfg *Config) error {
	var err error

	// 1. 加载配置
	viperLoader := viper.NewLoader(".", "config", "toml")
	viperLoader.SetDefault(SrvHttpAddr, ":8080")
	viperLoader.SetDefault(SrvGrpcAddr, ":8081")
	viperLoader.SetDefault(SrvSwaggerOpen, false)
	if err = config.AddSource(ConfigSourceFile, viperLoader); err != nil {
		//panic(errors.Errorf("加载配置错误:%s", err))
		return errors.Errorf("加载配置错误:%s", err)
	}
	for _, cc := range cfg.ConfigSource {
		if err = config.AddSource(cc.Source, cc.Cfg); err != nil {
			//panic(errors.Errorf("加载配置错误:%s", err))
			return errors.Errorf("加载配置错误:%s", err)
		}
	}

	// 2. 设置logger
	if err = app.setLogger(cfg.Logger); err != nil {
		return errors.Errorf("设置日志错误:%s", err)
	}

	// 3. 连接mysql/redis
	if err = app.setGorm(viperLoader); err != nil {
		return errors.Errorf("设置日志错误:%s", err)
	}
	app.setRedis(viperLoader)

	//4. http/grpc
	if vv := config.Get(ConfigSourceFile, SrvHttpAddr); vv != nil && reflect.TypeOf(vv).Kind() == reflect.String {
		var swaggerOpen bool
		if so := config.Get(ConfigSourceFile, SrvSwaggerOpen); so != nil && reflect.TypeOf(so).Kind() == reflect.Bool {
			swaggerOpen = so.(bool)
		}
		app.setHttpServer(vv.(string), cfg.Router, swaggerOpen)
	}
	if cfg.GrpcRegister != nil {
		if vv := config.Get(ConfigSourceFile, SrvGrpcAddr); vv != nil && reflect.TypeOf(vv).Kind() == reflect.String {
			app.setGrpcServer(vv.(string), cfg.GrpcRegister)
		}
	}

	return nil
}

func (app *App) setLogger(l log.Logger) error {
	if l != nil {
		log.SetLogger(l)
		return nil
	}
	logLevelStr, err := app.cmd.Flags().GetString(ArgLogLevel)
	if err != nil {
		return errors.Errorf("日志等级配置获取错误:%s", err)
	}
	logLevel, err := logrusLib.ParseLevel(logLevelStr)
	if err != nil {
		return errors.Errorf("日志等级配置错误:%s", err)
	}

	logFormatStr, err := app.cmd.Flags().GetString(ArgLogFormat)
	if err != nil {
		return errors.Errorf("日志格式配置获取错误:%s", err)
	}

	logPretty, err := app.cmd.Flags().GetBool(ArgLogPretty)
	if err != nil {
		return errors.Errorf("日志JSON美化输出配置获取错误:%s", err)
	}

	logCaller, err := app.cmd.Flags().GetBool(ArgLogCaller)
	if err != nil {
		return errors.Errorf("日志Caller配置获取错误:%s", err)
	}

	log.SetLogger(logrus.NewLogger(
		logrus.WithLevel(logLevel), logrus.WithFormatter(logFormatStr, logPretty), logrus.WithEnableCaller(logCaller),
	))
	return nil
}

func (app *App) setGorm(viperLoader *viper.Loader) error {
	mysqlCfg := viperLoader.Data.Sub("mysql")
	if mysqlCfg != nil {
		var conn *gormLib.DB
		var err error
		conn, err = gorm.New(&gorm.ConnConfig{
			Dsn:         mysqlCfg.GetString(MysqlDsn),
			MaxIdleConn: mysqlCfg.GetDuration(MysqlMaxIdleConn),
			MaxOpenConn: mysqlCfg.GetInt(MysqlMaxOpenConn),
			MaxLifeTime: mysqlCfg.GetDuration(MysqlMaxLifeTime),
			Logger: logger.New(
				log.NewHelper(log.Default()), // io writer
				logger.Config{
					SlowThreshold: time.Second, // 慢 SQL 阈值
					LogLevel:      logger.Info, // Log level
					Colorful:      true,        // 禁用彩色打印
				},
			),
		})
		if err != nil {
			//panic(errors.Errorf("初始化数据库连接错误:%s", err))
			return errors.Errorf("初始化数据库连接错误:%s", err)
		}
		app.DB = conn
	}
	return nil
}

func (app *App) setRedis(viperLoader *viper.Loader) {
	redisCfg := viperLoader.Data.Sub("redis")
	if redisCfg != nil {
		conn := redis.New(&redis.Cfg{
			Address:  redisCfg.GetString(RedisAddr),
			Password: redisCfg.GetString(RedisPassword),
			DB:       redisCfg.GetInt(RedisDb),
			PoolSize: redisCfg.GetInt(RedisPoolsize),
		})
		app.Redis = conn
	}
}

func (app *App) setHttpServer(addr string, router func(router gin.IRouter), swaggerOpen bool) {
	app.http = http_svc.NewServer(&http_svc.Config{
		Addr:        addr,
		Router:      router,
		SwaggerOpen: swaggerOpen,
	})
}

func (app *App) setGrpcServer(addr string, register func(srv *grpc.Server)) {
	srv := grpc.NewServer()
	register(srv)
	app.grpc = srv
	app.grpcAddr = addr
}

func (app *App) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				if app.http != nil {
					if err := app.http.Shutdown(ctx); err != nil {
						return err
					}
				}
				if app.grpc != nil {
					app.grpc.GracefulStop()
				}
				return ctx.Err()
			case <-ch:
				cancel()
			}
		}
	})
	if app.http != nil {
		eg.Go(app.http.ListenAndServe)
	}
	if app.grpc != nil {
		eg.Go(func() error {
			lis, err := net.Listen("tcp", app.grpcAddr)
			if err != nil {
				return err
			}
			return app.grpc.Serve(lis)
		})
	}
	if err := eg.Wait(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
		return err
	}

	return nil
}

func (app *App) Execute() error {
	return app.cmd.Execute()
}

func (app *App) stop() {
	var err error
	if app.Redis != nil {
		if err = app.Redis.Close(); err != nil {
			log.Error("关闭Redis错误", err)
		}
	}
}

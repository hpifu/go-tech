package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hpifu/go-account/pkg/account"
	godtoken "github.com/hpifu/go-godtoken/api"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/hpifu/go-kit/logger"
	"github.com/hpifu/go-tech/internal/es"
	"github.com/hpifu/go-tech/internal/mysql"
	"github.com/hpifu/go-tech/internal/service"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

// AppVersion name
var AppVersion = "unknown"

func main() {
	version := flag.Bool("v", false, "print current version")
	configfile := flag.String("c", "configs/tech.json", "config file path")
	flag.Parse()
	if *version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}

	// load config
	config := viper.New()
	config.SetEnvPrefix("tech")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()
	config.SetConfigType("json")
	fp, err := os.Open(*configfile)
	if err != nil {
		panic(err)
	}
	err = config.ReadConfig(fp)
	if err != nil {
		panic(err)
	}

	// init logger
	infoLog, warnLog, accessLog, err := logger.NewLoggerGroupWithViper(config.Sub("logger"))
	if err != nil {
		panic(err)
	}

	// init mysql
	db, err := mysql.NewMysql(config.GetString("mysql.uri"))
	if err != nil {
		panic(err)
	}
	infoLog.Infof("init mysql success. uri [%v]", config.GetString("mysql.uri"))

	// init elasticsearch
	esclient, err := es.NewES(
		config.GetString("es.uri"),
		config.GetString("es.index"),
		config.GetDuration("es.timeout"),
	)
	if err != nil {
		panic(err)
	}
	infoLog.Infof("init elasticsearch success. uri [%v]", config.GetString("es.uri"))

	// init http client
	accountCli := account.NewClient(
		config.GetString("account.address"),
		config.GetInt("account.maxConn"),
		config.GetDuration("account.connTimeout"),
		config.GetDuration("account.recvTimeout"),
	)

	// init godtoken client
	var kacp = keepalive.ClientParameters{
		Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
		Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
		PermitWithoutStream: true,             // send pings even without active streams
	}
	conn, err := grpc.Dial(
		config.GetString("godtoken.address"),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		panic(err)
	}
	godtokenCli := godtoken.NewServiceClient(conn)
	infoLog.Infof("init godtoken client success. address: [%v]", config.GetString("godtoken.address"))

	// init services
	svc := service.NewService(db, esclient, accountCli, godtokenCli)
	svc.SetLogger(infoLog, warnLog, accessLog)

	// init gin
	origins := config.GetStringSlice("service.allowOrigins")
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowMethods:     []string{"PUT", "POST", "GET", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "Accept", "Cache-Control", "X-Requested-With"},
		AllowCredentials: true,
	}))

	// set handler
	d := hhttp.NewGinHttpDecorator(infoLog, warnLog, accessLog)
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	r.GET("/article", d.Decorate(svc.GETArticles))
	r.GET("/article/:id", d.Decorate(svc.GETArticle))
	r.GET("/articles/author/:authorID", d.Decorate(svc.GETArticlesAuthor))
	r.GET("/articles/tag/:tag", d.Decorate(svc.GETArticlesTag))
	r.POST("/article", d.Decorate(svc.POSTArticle))
	r.PUT("/article/:id", d.Decorate(svc.PUTArticle))
	r.DELETE("article/:id", d.Decorate(svc.DELETEArticle))
	r.GET("/tagcloud", d.Decorate(svc.GETTagCloud))
	r.GET("/search", d.Decorate(svc.Search))
	r.POST("/like/:id", d.Decorate(svc.Like))

	infoLog.Infof("%v init success, setting: %v", os.Args[0], config.AllSettings())

	// run server
	server := &http.Server{
		Addr:    config.GetString("service.port"),
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// graceful quit
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	infoLog.Infof("%v shutdown ...", os.Args[0])

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		warnLog.Errorf("%v shutdown fail or timeout", os.Args[0])
		return
	}
	warnLog.Out.(*rotatelogs.RotateLogs).Close()
	accessLog.Out.(*rotatelogs.RotateLogs).Close()
	infoLog.Errorf("%v shutdown success", os.Args[0])
	infoLog.Out.(*rotatelogs.RotateLogs).Close()
}

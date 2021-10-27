package main

import (
	"errors"
	"net/http"
	"time"

	"apiserver/config"
	"apiserver/router"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	runmode      = "runmode"
	addr         = "addr"
	url          = "url"
	maxPingCount = "max_ping_count"
)

var (
	// cfg 变量值由命令行 flag 传入
	// e.g.可以传值(./apiserver -c config.yaml)
	cfg = pflag.StringP("config", "c", "", "apiserver config file path.")
)

func main() {
	pflag.Parse()

	// init config
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// Set gin mode.
	// gin 有 3 种运行模式：debug\release\test, 其中 debug 模式会打印很多 debug 信息
	gin.SetMode(viper.GetString(runmode))

	// Create the Gin engine.
	g := gin.New()

	var middlewares []gin.HandlerFunc // func(*Context)

	// Routes.
	router.Load(
		// Cores.
		g,

		// Middlwares.
		middlewares...,
	)

	// Ping the server to make sure the router is working.
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()
	log.Infof("Start to listening the incoming requests on http address: %s", viper.GetString(addr))
	log.Info(http.ListenAndServe(viper.GetString(addr), g).Error())
}

// pingServer pings the http server to make sure the router is working.
func pingServer() error {
	for i := 0; i < viper.GetInt(maxPingCount); i++ {
		// Ping the server by sending a GET request to `/health`.
		resp, err := http.Get(viper.GetString(url) + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}

		// Sleep for a second to continue the next ping.
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	//goland:noinspection GoErrorStringFormat
	return errors.New("Cannot connect to the router.")
}

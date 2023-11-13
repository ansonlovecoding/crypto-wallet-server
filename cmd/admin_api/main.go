package main

import (
	admin "Share-Wallet/internal/admin_api"
	"Share-Wallet/pkg/common/config"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	_ "Share-Wallet/pkg/swagger/admin_api"
	"Share-Wallet/pkg/utils"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strconv"
)

// @title           Admin API
// @version         1.0.0
// @description     This is server for management platform.
// @license.name  Apache 2.0

// @host      cms.wallet.com
// @BasePath  /cms/v1
func main() {
	//runtime.SetMutexProfileFraction(1) // 开启对锁调用的跟踪
	//runtime.SetBlockProfileRate(1)     // 开启对阻塞操作的跟踪
	//
	//go func() {
	//	// 启动一个自定义mux的http服务器
	//	mux := http.NewServeMux()
	//	mux.HandleFunc("/debug/pprof/", pprof.Index)
	//	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	//	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	//	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	//	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	//
	//	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
	//		w.Write([]byte("hello"))
	//	})
	//	// 启动一个 http server，注意 pprof 相关的 handler 已经自动注册过了
	//	if err := http.ListenAndServe(":6062", mux); err != nil {
	//		log1.Fatal(err)
	//		log.NewError("", utils.GetSelfFuncName(), "启动pprof报错：", err.Error())
	//	}
	//	os.Exit(0)
	//}()

	log.NewPrivateLog(constant.AdminAPILogFileName)
	gin.SetMode(gin.DebugMode)
	router := admin.NewGinRouter()
	router.Use(utils.CorsHandler())
	defaultPorts := config.Config.AdminApi.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10006 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.AdminApi.ListenIP != "" {
		address = config.Config.AdminApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	address = config.Config.AdminApi.ListenIP + ":" + strconv.Itoa(*ginPort)
	fmt.Println("start cms api server, address: ", address)
	if config.Config.SwaggerEnable {
		//http://cms.wallet.com/swagger/index.html
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
	router.Run(address)

}

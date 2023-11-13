package main

import (
	"Share-Wallet/internal/cronjob/logic"
	"Share-Wallet/pkg/common/constant"
	"Share-Wallet/pkg/common/log"
	_ "Share-Wallet/pkg/swagger/admin_api"
	"fmt"
	"sync"
)

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

	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.CronLogFileName)
	fmt.Println("start cronjob server")
	logic.Run()
	wg.Wait()

}

package router

import (
	"context"
	"fmt"
	"keyboardComment/dbs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func SetupRouter() {
	Router = gin.Default()

	ManagerRouter("manager")

	Router.NoRoute(func(ctx *gin.Context) {
		ctx.String(http.StatusNotFound, "404")
	})

	err := dbs.InitEnvironment("pro")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer dbs.Close()

	server := &http.Server{
		Addr:           ":8088",
		Handler:        Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   100 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		err = server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// 优雅退出
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	// 卡壳，接受控制台退出命令继续往下执行
	<-ch
	cxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Shutdown(cxt)
	if err != nil {
		log.Println("err", err)
	}
}

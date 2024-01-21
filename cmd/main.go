package main

import (
	"github.com/gin-contrib/static"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slot-crawler/internal/database"
	"slot-crawler/internal/server/router"
	"syscall"
	"time"
)

func main() {
	//
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}

	database.Initialize()

	r := router.InitRouter()
	r.Use(static.Serve("/", static.LocalFile("./web", false)))
	r.ForwardedByClientIP = true
	err := r.SetTrustedProxies([]string{"127.0.0.1"})
	if err != nil {
		log.Panic("set trusted proxies fail")
		return
	}
	srv := &http.Server{
		Addr:        ":8081",
		Handler:     r,
		ReadTimeout: 10 * time.Second,
		//WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err.Error())
		return
	}
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be caught, so don't need added it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		database.DisConnect()
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}

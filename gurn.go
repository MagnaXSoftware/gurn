package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"magnax.ca/gurn/pkg/config"
	"magnax.ca/gurn/pkg/storage"
	"magnax.ca/gurn/pkg/web"
)

func main() {
	pConfigFile := flag.String("config", "gurn.hcl", "Path to the config file")
	flag.Parse()

	rawConf, err := os.ReadFile(*pConfigFile)
	if err != nil {
		fmt.Printf("Error reading config file: %s", err.Error())
		return
	}

	conf, err := config.ParseConfig(rawConf)
	if err != nil {
		fmt.Printf("Error parsing config file: %s", err.Error())
		return
	}

	err = storage.InitDatabases(conf.Database)
	if err != nil {
		fmt.Printf("Error during db initialization: %s", err.Error())
		return
	}

	if _, ok := os.LookupEnv("GIN_MODE"); !ok {
		gin.SetMode(gin.DebugMode)
	}

	s := web.NewServer(conf)
	go func() {
		err = s.Run()
		if err != nil {
			fmt.Printf("Error during run: %s", err.Error())
			return
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")

}

package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"

	"magnax.ca/gurn/pkg/storage"
	"magnax.ca/gurn/pkg/web"
)

func main() {
	err := storage.InitDatabases()
	if err != nil {
		fmt.Printf("Error during db initialization: %s", err.Error())
		return
	}

	if _, ok := os.LookupEnv("GIN_MODE"); !ok {
		gin.SetMode(gin.DebugMode)
	}

	s := web.NewServer().WithAddr("localhost:8080")
	err = s.Run()
	if err != nil {
		fmt.Printf("Error during run: %s", err.Error())
		return
	}
}

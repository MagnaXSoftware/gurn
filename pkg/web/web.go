package web

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"magnax.ca/gurn/pkg/config"
)

type Server struct {
	Addr string

	e   *gin.Engine
	srv *http.Server
}

func NewServer(conf *config.Config) *Server {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"})

	MountFuncs(r.Group("/_"))

	r.GET("/:urn", FindAndRedirectByURN)
	r.GET("/", IndexPage)
	r.GET("/favicon.ico", func(c *gin.Context) { c.AbortWithStatus(http.StatusNotFound) })

	s := &Server{e: r}
	s.Addr = conf.BindAddr

	return s
}

func (s *Server) Run() error {
	s.srv = &http.Server{
		Addr:    s.Addr,
		Handler: s.e,
	}

	err := s.srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.srv == nil {
		return fmt.Errorf("webserver was never started")
	}
	return s.srv.Shutdown(ctx)
}

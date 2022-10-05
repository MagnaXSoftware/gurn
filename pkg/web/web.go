package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	*gin.Engine

	Addr string
}

func NewServer() *Server {
	r := gin.Default()
	_ = r.SetTrustedProxies([]string{"192.168.0.0/16", "172.16.0.0/12", "10.0.0.0/8"})

	MountFuncs(r.Group("/_"))

	r.GET("/:urn", FindAndRedirectByURN)
	r.GET("/", IndexPage)
	r.GET("/favicon.ico", func(c *gin.Context) { c.AbortWithStatus(http.StatusNotFound) })

	return &Server{Engine: r}
}

func (s *Server) WithAddr(Addr string) *Server {
	s.Addr = Addr
	return s
}

func (s *Server) Run() error {
	return s.Engine.Run(s.Addr)
}

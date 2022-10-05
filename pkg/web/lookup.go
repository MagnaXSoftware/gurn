package web

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"

	"magnax.ca/gurn/pkg/storage"
)

type Result string

const (
	NotFound Result = "notfound"
	Success  Result = "success"
)

func FindAndRedirectByURN(c *gin.Context) {
	urnKey, err := url.QueryUnescape(c.Param("urn"))
	if err != nil {
		// TODO display list of similar or all known urns
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	urn := storage.Urn{}.FindByName(urnKey)
	if urn == nil {
		go logAccess(urnKey, NotFound)
		c.String(http.StatusNotFound, "Was not able to find URN matching name %s", urnKey)
		return
	}

	go logAccess(urnKey, Success)

	c.Redirect(http.StatusFound, urn.Destination)
}

func logAccess(key string, result Result) {
	a := &storage.Access{Urn: key, AccessedAt: time.Now(), Result: string(result)}
	_ = a.Save()
}

func IndexPage(c *gin.Context) {
	c.String(http.StatusOK, "Listing everything")
}

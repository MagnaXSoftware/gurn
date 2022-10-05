package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"magnax.ca/gurn/pkg/storage"
)

func MountFuncs(group *gin.RouterGroup) {
	group.GET("/stats/:urn", StatsByURN)
	group.POST("/add")
	// group.GET("/add")
}

func StatsByURN(c *gin.Context) {
	urnKey, err := url.QueryUnescape(c.Param("urn"))
	if err != nil {
		_ = c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	accesses := storage.Access{}.FindByUrn(urnKey)
	p, err := json.MarshalIndent(accesses, "\t", "\t")
	if err != nil {
		panic(err)
	}

	message := fmt.Sprintf("Looking up stats for %s\n\nThere are %d entries for this key: \n\t%v", urnKey, len(accesses), string(p))
	c.String(http.StatusOK, message)
}

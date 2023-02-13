package decision

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/paging"
	"io.github.aliaksandrrachko/decision-helper/webapp/pkg/response"
)

// todo: see logger and authHandler !!!

func RegisterHandlers(router gin.IRouter, service Service, logger logrus.Logger) {
	res := resource{service, logger}

	router.GET("/decisions/:id", res.get)
	router.GET("/decisions", res.query)
	router.POST("decisions", res.create)
	router.PUT("/decisions/:id", res.update)
	router.DELETE("/decisions/:id", res.delete)
}

type resource struct {
	service Service
	logger  logrus.Logger
}

func (r resource) get(c *gin.Context) {
	id, parseError := strconv.ParseInt(c.Param("id"), 0, 8)
	if parseError != nil {
		c.IndentedJSON(http.StatusBadRequest, response.NewErrorResponseByCode(http.StatusBadRequest, parseError))
		return
	}

	decision, err := r.service.Get(c.Request.Context(), id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, response.NewErrorResponseByCode(http.StatusNotFound, err))
		return
	}

	c.IndentedJSON(http.StatusOK, decision)
}

func (r resource) query(c *gin.Context) {
	metaInfo := paging.NewMetaInfoFromRequest(c.Request)
	decisions, err := r.service.Query(c.Request.Context(), metaInfo.Offset, metaInfo.Limit)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, response.NewErrorResponseByCode(http.StatusNotFound, err))
		return
	}

	c.IndentedJSON(http.StatusOK, paging.New(decisions, metaInfo))
}

func (r resource) create(c *gin.Context) {
	var decision Decision

	if err := c.BindJSON(&decision); err != nil {
		c.IndentedJSON(http.StatusBadRequest, response.NewErrorResponseByCode(http.StatusBadRequest, err))
		return
	}

	id, err := r.service.Create(c.Request.Context(), decision)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, response.NewErrorResponseByCode(http.StatusInternalServerError, err))
		return
	}

	c.Header("Location", fmt.Sprintf("%s/%d", (c.Request.Host+c.Request.URL.Path), id))
	c.Status(http.StatusCreated)
}

func (r resource) update(c *gin.Context) {
	id, parseError := strconv.ParseInt(c.Param("id"), 0, 64)
	if parseError != nil {
		c.IndentedJSON(http.StatusBadRequest, response.NewErrorResponseByCode(http.StatusBadRequest, parseError))
		return
	}

	var decision Decision

	if err := c.BindJSON(&decision); err != nil {
		c.IndentedJSON(http.StatusBadRequest, response.NewErrorResponseByCode(http.StatusBadRequest, err))
		return
	}

	err := r.service.Update(c.Request.Context(), id, decision)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, response.NewErrorResponseByCode(http.StatusInternalServerError, err))
		return
	}

	c.Status(http.StatusOK)
}

func (r resource) delete(c *gin.Context) {
	id, parseError := strconv.ParseInt(c.Param("id"), 0, 64)
	if parseError != nil {
		c.IndentedJSON(http.StatusBadRequest, response.NewErrorResponseByCode(http.StatusBadRequest, parseError))
	}

	r.service.Delete(c, id)
	c.Status(http.StatusNoContent)
}

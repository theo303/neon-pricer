package http

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/usecases"

	"github.com/gin-gonic/gin"
)

const (
	siliconeParam = "silic"
)

type configHandlers struct {
	config *conf.Configuration
}

func (ch configHandlers) getConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "config.html", ch.config)
	}
}

func parseParamInPostForm(c *gin.Context, key string) (float64, error) {
	valueStr := c.PostForm(key)
	if valueStr == "" {
		return 0, fmt.Errorf("key %s not found", key)
	}
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0, fmt.Errorf("parsing %s: %w", valueStr, err)
	}
	return value, nil
}

func (ch configHandlers) setConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Printf("setConfig: error while reading body: %s", err)
			c.Status(http.StatusBadRequest)
			return
		}

		ch.config, err = usecases.UpdateConfigWithPostForm(ch.config, body)
		if err != nil {
			fmt.Printf("setConfig: error while updating config with body %s: %s", string(body), err)
			c.Status(http.StatusBadRequest)
			return
		}

		c.Status(http.StatusNoContent)
	}
}

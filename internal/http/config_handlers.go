package http

import (
	"fmt"
	"io"
	"net/http"

	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/usecases"

	"github.com/gin-gonic/gin"
)

type configHandlers struct {
	config *conf.Configuration
}

func (ch configHandlers) getConfig() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "config.html", ch.config)
	}
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

type radioButton struct {
	Name      string
	IsDefault bool
}

func (ch configHandlers) getInput() gin.HandlerFunc {
	return func(c *gin.Context) {
		data := struct {
			Plexis []radioButton
		}{}
		for _, plexi := range ch.config.Plexis {
			data.Plexis = append(data.Plexis, radioButton{
				Name:      plexi.Name,
				IsDefault: plexi.Name == "incolore",
			})
		}
		c.HTML(http.StatusOK, "input.html", data)
	}
}

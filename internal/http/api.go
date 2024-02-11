package http

import (
	"fmt"
	"math"
	"net/http"

	"theo303/neon-pricer/conf"
	"theo303/neon-pricer/internal/svg"
	"theo303/neon-pricer/internal/usecases"

	"github.com/gin-gonic/gin"
)

type API struct {
	config *conf.Configuration
	port   int

	configHandlers configHandlers
}

func NewAPI(config conf.Configuration, port int) API {
	return API{
		config: &config,
		port:   port,
		configHandlers: configHandlers{
			config: &config,
		},
	}
}

func (a API) Run() error {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/config", a.configHandlers.getConfig())
	r.POST("/config", a.configHandlers.setConfig())
	r.GET("/input", a.configHandlers.getInput())
	r.POST("/compute", a.compute())

	return r.Run(fmt.Sprintf(":%d", a.port))
}

type computationResult struct {
	Group         string
	LengthPx      float64
	LengthMm      float64
	WidthMm       float64
	HeightMm      float64
	SiliconePrice float64
	LedPrice      float64
	PlexiPrice    float64
}
type resultData struct {
	Results []computationResult
}

func (a API) compute() gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			_ = c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		fileBuf, err := file.Open()
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		forms, err := svg.RetrieveForms(fileBuf, "")
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		sizes, err := usecases.GetSizes(forms, a.config.Scale)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		fmt.Println(c.PostForm("plexi"))
		prices, err := usecases.GetPrice(a.config.Pricing, sizes, c.PostForm("plexi"))
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var resData resultData
		for g, size := range sizes {
			resData.Results = append(resData.Results, computationResult{
				Group:         g,
				LengthPx:      math.Round(size.LengthPx),
				LengthMm:      math.Round(size.Length),
				WidthMm:       math.Round(size.Width),
				HeightMm:      math.Round(size.Height),
				SiliconePrice: prices[g].SiliconePrice,
				LedPrice:      prices[g].LEDPrice,
				PlexiPrice:    prices[g].PlexiPrice,
			})
		}

		c.HTML(http.StatusOK, "response.html", resData)
	}
}

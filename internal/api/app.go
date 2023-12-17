package api

import (
	"fmt"
	"math"
	"net/http"
	"theo303/neon-pricer/configuration"
	"theo303/neon-pricer/internal/svg"
	"theo303/neon-pricer/internal/usecases"

	"github.com/gin-gonic/gin"
)

type Conf struct {
	configuration.Configuration
	Port int
}

func Run(conf Conf) error {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/compute", compute(conf.Scale))

	return r.Run(fmt.Sprintf(":%d", conf.Port))
}

type computationResult struct {
	Group    string
	LengthPx float64
	LengthMm float64
	WidthMm  float64
	HeightMm float64
}
type resultData struct {
	Results []computationResult
}

func compute(scale float64) gin.HandlerFunc {
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

		lengths, err := usecases.GetLengths(forms)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		bounds, err := usecases.GetBounds(forms)
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		var resData resultData
		for g := range forms {
			fmt.Println(lengths[g] * 1000.0 / scale)
			resData.Results = append(resData.Results, computationResult{
				Group:    g,
				LengthPx: math.Round(lengths[g]),
				LengthMm: math.Round(lengths[g] * 1000.0 / scale),
				WidthMm:  math.Round(bounds[g].Width() * 1000.0 / scale),
				HeightMm: math.Round(bounds[g].Height() * 1000.0 / scale),
			})

			fmt.Println(resData)
		}

		c.HTML(http.StatusOK, "response.html", resData)
	}
}

package gin_middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger"
	"github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/dto/response"
)

func LogRequest(logger logger.Logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		method := c.Request.Method
		url := c.Request.URL.String()
		headers := c.Request.Header

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.ErrorResponse{
				Error: "something went wrong",
			})
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		logger.Infof("Method: %s\nURL: %s\nHeaders: %v\nBody: %s", method, url, headers, string(bodyBytes))

		c.Next()
	}
}

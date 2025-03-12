package response

import (
	"github.com/gin-gonic/gin"
	"github.com/cryptoPickle/go-ddd-example/internal/common/errors"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WithHTTPError(err error, c *gin.Context) {
	httpErr := errors.ToHTTPError(err)
	c.JSON(httpErr.StatusCode, ErrorResponse{
		Error: httpErr.Reason,
	})
}

package restapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cryptoPickle/go-ddd-example/internal/application/interfaces"
	"github.com/cryptoPickle/go-ddd-example/internal/infrastructure/logger"
	restapi_mapper "github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/dto/mapper"
	"github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/dto/request"
	"github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/dto/response"
	gin_middleware "github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/middleware"
)

type PayoutController struct {
	service interfaces.PayoutService
	logger  logger.Logger
}

func NewPayoutController(g *gin.Engine, service interfaces.PayoutService, logger logger.Logger) *PayoutController {
	controller := PayoutController{
		service: service,
	}

	g.Use(gin_middleware.LogRequest(logger))
	g.POST("/payout", controller.CreatePayouts)
	return &controller
}

// @Summary Creates Payouts for Sellers
// @Tags Payouts
// @Description Takes desired currency and multiple item ids, generates payouts respecting the transaction limit, if limit is exceed generates multiple payouts
// @Accept json
// @Produce json
// @Param payout body request.CreatePayoutRequest true "Payout Details"
// @Success 201 {object} response.CreatePayoutResponse "Success Response"
// @Failure 400 {object} response.ErrorResponse "Bad Request"
// @Failure 500 {object} response.ErrorResponse "Internal Server Error"
// @Router /payout [post]
func (pc *PayoutController) CreatePayouts(c *gin.Context) {
	var createPayoutRequest request.CreatePayoutRequest
	if err := c.Bind(&createPayoutRequest); err != nil {
		response.WithHTTPError(err, c)
		return
	}

	payoutCommand, err := createPayoutRequest.ToCreatePayoutCommand()
	if err != nil {
		response.WithHTTPError(err, c)
		return
	}

	commandResult, err := pc.service.CreatePayouts(c, payoutCommand)
	if err != nil {
		response.WithHTTPError(err, c)
		return
	}

	response := restapi_mapper.ToPayoutResponse(commandResult)

	c.JSON(http.StatusCreated, response)
}

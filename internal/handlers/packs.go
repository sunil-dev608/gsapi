package handlers

import (
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type PacksHandler struct {
	packSizes []int
}

// NewTransactionHandler returns a new instance of TransactionHandler
func NewPacksHandler(packSizes []int) *PacksHandler {
	return &PacksHandler{
		packSizes: packSizes,
	}
}

func (ph *PacksHandler) getPacksForItems(numberOfItems int) map[int]int {

	packs := make(map[int]int)
	for _, packSize := range ph.packSizes {
		for numberOfItems >= packSize {
			packs[packSize]++
			numberOfItems -= packSize
		}
	}
	if numberOfItems > 0 {
		packs[ph.packSizes[len(ph.packSizes)-1]]++
	}
	return packs
}

func (ph *PacksHandler) GetPacksForItems(c echo.Context) error {

	logger := ctxzap.Extract(c.Request().Context())
	var req struct {
		Items int `json:"items"`
	}

	if err := c.Bind(&req); err != nil {
		logger.Warn("bad request", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "bad request"})
	}

	if req.Items <= 0 {
		logger.Warn("bad request", zap.Int("items", req.Items))
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "items must be greater than 0"})
	}

	response := ph.getPacksForItems(req.Items)
	c.JSON(http.StatusOK, response)
	return nil
}

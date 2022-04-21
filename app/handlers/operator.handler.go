package handlers

import (
	services "database-ms/app/services"

	"github.com/gin-gonic/gin"
)

type OperatorHandler struct {
	service services.OperatorServiceInterface
}

func NewOperatorAPI(operatorService services.OperatorServiceInterface) *OperatorHandler {
	return &OperatorHandler{service: operatorService}
}

func (handler *OperatorHandler) CreateOperator(ctx *gin.Context) {
	// TODO
}

// Get operators for the entire org
func (handler *OperatorHandler) GetOperators(ctx *gin.Context) {
	// TODO
}

func (handler *OperatorHandler) UpdateOperator(ctx *gin.Context) {
	// TODO
}

func (handler *OperatorHandler) DeleteOperator(ctx *gin.Context) {
	// TODO
}
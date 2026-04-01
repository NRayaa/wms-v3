package controller

import (
	"net/http"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
)

// BuyerController defines handlers for buyer resources.
type BuyerController struct {
	service services.BuyerService
}

// NewBuyerController constructor.
func NewBuyerController(service services.BuyerService) *BuyerController {
	return &BuyerController{service: service}
}

// CreateBuyer endpoint.
func (ctrl *BuyerController) CreateBuyer(c *gin.Context) {
	var payload services.CreateBuyerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}

	buyer, err := ctrl.service.CreateBuyer(payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}

	utils.SendSuccess(c, buyer, "Buyer berhasil ditambahkan", http.StatusCreated)
}

// GetBuyerByID endpoint.
func (ctrl *BuyerController) GetBuyerByID(c *gin.Context) {
	id := c.Param("id")
	buyer, err := ctrl.service.GetBuyerByID(id)
	if err != nil {
		utils.SendError(c, 404, err.Error())
		return
	}
	utils.SendSuccess(c, buyer, "", http.StatusOK)
}

// ListBuyers endpoint.
func (ctrl *BuyerController) ListBuyers(c *gin.Context) {
	buyers, err := ctrl.service.ListBuyers()
	if err != nil {
		utils.SendError(c, 500, err.Error())
		return
	}
	utils.SendSuccess(c, buyers, "", http.StatusOK)
}

// UpdateBuyer endpoint.
func (ctrl *BuyerController) UpdateBuyer(c *gin.Context) {
	id := c.Param("id")
	var payload services.UpdateBuyerPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		validationErrors := []utils.ErrorItem{{Field: "", Message: err.Error()}}
		utils.SendValidationError(c, validationErrors)
		return
	}
	buyer, err := ctrl.service.UpdateBuyer(id, payload)
	if err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, buyer, "Buyer berhasil diupdate", http.StatusOK)
}

// DeleteBuyer endpoint.
func (ctrl *BuyerController) DeleteBuyer(c *gin.Context) {
	id := c.Param("id")
	if err := ctrl.service.DeleteBuyer(id); err != nil {
		utils.SendError(c, 400, err.Error())
		return
	}
	utils.SendSuccess(c, nil, "Buyer berhasil dihapus", http.StatusOK)
}

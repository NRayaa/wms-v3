// ListAllProductMastersHandler menampilkan seluruh data master secara descending
package controller

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"wms/models"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type InboundRequest struct {
	Name       string  `json:"name" binding:"required"`
	Item       int     `json:"item" binding:"required,gt=0"`
	Price      float64 `json:"price" binding:"required"`
	CategoryID *string `json:"category_id,omitempty"`
	StickerID  *string `json:"sticker_id,omitempty"`
}

func generateUniqueBarcode() string {
	t := time.Now().UnixNano()
	r := rand.Intn(100000)
	return fmt.Sprintf("BC-%d-%d", t, r)
}

func ListAllProductMastersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var masters []models.ProductMaster
		if err := db.Order("timestamp DESC").Find(&masters).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendSuccess(c, masters, "List master data", http.StatusOK)
	}
}

func InboundManualHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req InboundRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			var verrs []utils.ErrorItem
			if ve, ok := err.(validator.ValidationErrors); ok {
				for _, ferr := range ve {
					verrs = append(verrs, utils.ErrorItem{
						Field:   ferr.Field(),
						Message: ferr.Error(),
					})
				}
			} else {
				verrs = append(verrs, utils.ErrorItem{Field: "", Message: err.Error()})
			}
			utils.SendValidationError(c, verrs)
			return
		}

		barcode := generateUniqueBarcode()
		barcodeWarehouse := generateUniqueBarcode()

		// Logic BE: tentukan category_id/sticker_id otomatis

		var categoryID, stickerID, typeID string
		if req.Price >= 100000 {
			if req.CategoryID != nil {
				categoryID = *req.CategoryID
			}
			stickerID = ""
			typeID = "categories"
		} else {
			if req.StickerID != nil {
				stickerID = *req.StickerID
			}
			categoryID = ""
			typeID = "sticker"
		}

		master := models.ProductMaster{
			Barcode:          barcode,
			BarcodeWarehouse: barcodeWarehouse,
			Name:             req.Name,
			Item:             req.Item,
			Price:            req.Price,
			CategoryID:       categoryID,
			StickerID:        stickerID,
			TypeID:           typeID,
			Location:         "rack",
			TypeOut:          "cargo",
		}

		// Insert ke ProductPending
		pending := models.ProductPending{
			Barcode: barcode,
			Name:    req.Name,
			Item:    req.Item,
			Price:   req.Price,
			Status:  "non", // default status valid
		}
		if err := db.Create(&pending).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		if err := db.Create(&master).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}

		utils.SendSuccess(c, gin.H{"pending": pending, "master": master}, "Inbound berhasil dibuat", http.StatusOK)
	}
}

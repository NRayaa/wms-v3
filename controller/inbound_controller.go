// ListAllProductMastersHandler menampilkan seluruh data master secara descending
package controller

import (
	"net/http"

	"wms/models"
	"wms/services"
	"wms/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ListAllProductMastersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var masters []models.ProductMaster
		if err := db.Order("created_at DESC").Find(&masters).Error; err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendSuccess(c, masters, "List master data", nil, http.StatusOK)
	}
}

// Tambahkan variabel global untuk service
var inboundService = services.NewInboundService()

// Handler untuk upload bulk file (step 1)
func InboundBulkUploadHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		supplier := c.PostForm("supplier")
		typeProduct := c.PostForm("type_product") // reguler/sticker
		fileType := c.PostForm("type")            // csv/xlsx/xls

		file, header, err := c.Request.FormFile("file")
		if err != nil {
			utils.SendError(c, 400, "File tidak ditemukan")
			return
		}
		defer file.Close()

		// Parse file
		headers, rows, err := utils.ParseBulkFile(file, fileType)
		if err != nil {
			utils.SendError(c, 400, "Gagal membaca file: "+err.Error())
			return
		}

		// Kirim headers dan rows ke FE untuk proses mapping
		utils.SendSuccess(c, gin.H{
			"headers":      headers,
			"rows":         rows,
			"filename":     header.Filename,
			"supplier":     supplier,
			"type_product": typeProduct,
			"type":         fileType,
		}, "File berhasil diupload", nil, http.StatusOK)
	}
}

// Handler untuk proses bulk setelah mapping (step 2)
func InboundBulkProcessHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.BulkInboundRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.SendError(c, 400, "Request tidak valid: "+err.Error())
			return
		}

		inserted, skipped, skipDetails := inboundService.InboundBulkProcess(req, db)
		utils.SendSuccess(c, gin.H{
			"inserted":     inserted,
			"skipped":      skipped,
			"skip_details": skipDetails,
		}, "Bulk inbound selesai", nil, http.StatusOK)
	}
}

func InboundManualHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.InboundRequest
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

		pending, master, err := inboundService.InboundManual(models.InboundRequest{
			Name:       req.Name,
			Item:       req.Item,
			Price:      req.Price,
			CategoryID: req.CategoryID,
			StickerID:  req.StickerID,
			Status:     req.Status,
		}, db)
		if err != nil {
			utils.SendError(c, 500, err.Error())
			return
		}
		utils.SendSuccess(c, gin.H{"pending": pending, "master": master}, "Inbound berhasil dibuat", nil, http.StatusOK)
	}
}

package utils

import (
	"encoding/csv"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"
)

// BulkFileType represents supported file types for bulk upload
const (
	BulkFileTypeCSV  = "csv"
	BulkFileTypeXLSX = "xlsx"
	BulkFileTypeXLS  = "xls"
)

// ParseBulkFile parses uploaded file and returns headers and rows
func ParseBulkFile(file multipart.File, fileType string) (headers []string, rows [][]string, err error) {
	switch strings.ToLower(fileType) {
	case BulkFileTypeCSV:
		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, nil, err
			}
			if headers == nil {
				headers = record
			} else {
				rows = append(rows, record)
			}
		}
		return headers, rows, nil
	case BulkFileTypeXLSX, BulkFileTypeXLS:
		xlFile, err := excelize.OpenReader(file)
		if err != nil {
			return nil, nil, err
		}
		sheetName := xlFile.GetSheetName(0)
		if sheetName == "" {
			return nil, nil, fmt.Errorf("no sheet found")
		}
		rowsIter, err := xlFile.Rows(sheetName)
		if err != nil {
			return nil, nil, err
		}
		rowIdx := 0
		for rowsIter.Next() {
			row, _ := rowsIter.Columns()
			if rowIdx == 0 {
				headers = row
			} else {
				rows = append(rows, row)
			}
			rowIdx++
		}
		return headers, rows, nil
	default:
		return nil, nil, fmt.Errorf("unsupported file type: %s", fileType)
	}
}

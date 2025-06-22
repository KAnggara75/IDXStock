/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/gpl-3.0.html>.
 *
 * @author KAnggara75 on Sun 22/06/25 15.45
 * @project IDXStock handler
 * https://github.com/KAnggara75/IDXStock/tree/main/handler
 */

package handler

import (
	"fmt"
	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

func UploadExcel(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	filename := file.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid File Type"})
		return
	}

	prefix := config.GetStockListPrefix()
	if !strings.HasPrefix(filename, prefix) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Filename"})
		return
	}

	datePattern := regexp.MustCompile(`(\d{8})\.xlsx$`)
	matches := datePattern.FindStringSubmatch(filename)
	if matches == nil || len(matches) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data Not Update"})
		return
	}

	tanggalFile := matches[1]
	today := time.Now().Format("20060102")
	if tanggalFile != today {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      fmt.Sprintf("Data Not Update: %s", today),
			"file_date":  tanggalFile,
			"today_date": today,
		})
		return
	}

	tmpFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp file"})
		return
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {

		}
	}(tmpFile.Name())

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer func() {
		if err := src.Close(); err != nil {
			fmt.Printf("Warning: failed to close file: %v\n", err)
		}
	}()

	_, err = io.Copy(tmpFile, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		return
	}
	defer func() {
		if err := tmpFile.Close(); err != nil {
			fmt.Printf("Warning: Failed to copy file: %v\n", err)
		}
	}()

	f, err := excelize.OpenFile(tmpFile.Name())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Warning: Invalid Excel file: %v\n", err)
		}
	}()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read rows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename": filename,
		"data":     rows,
	})
}

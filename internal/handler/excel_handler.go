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
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File harus .xlsx"})
		return
	}

	const prefix = "Daftar Saham"
	if !strings.HasPrefix(filename, prefix) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama file harus diawali dengan 'Daftar Saham'"})
		return
	}

	datePattern := regexp.MustCompile(`(\d{8})\.xlsx$`)
	matches := datePattern.FindStringSubmatch(filename)
	if matches == nil || len(matches) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama file harus diakhiri dengan tanggal format YYYYMMDD sebelum .xlsx"})
		return
	}
	tanggalFile := matches[1]
	today := time.Now().Format("20060102")
	if tanggalFile != today {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":      fmt.Sprintf("Tanggal di filename harus hari ini: %s", today),
			"file_date":  tanggalFile,
			"today_date": today,
		})
		return
	}

	// Simpan sementara file, proses selanjutnya seperti biasa
	tmpFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp file"})
		return
	}
	defer os.Remove(tmpFile.Name())

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
		return
	}
	defer src.Close()
	_, err = io.Copy(tmpFile, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to copy file"})
		return
	}
	tmpFile.Close()

	f, err := excelize.OpenFile(tmpFile.Name())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
		return
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read rows"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename":  filename,
		"sheet":     sheetName,
		"rows":      rows,
		"validated": true,
	})
}

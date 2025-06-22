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
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"os"
)

func UploadExcel(c *gin.Context) {
	// Ambil file dari form (key: "file")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		return
	}

	// Validasi ekstensi
	if ext := file.Filename[len(file.Filename)-5:]; ext != ".xlsx" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only .xlsx files allowed"})
		return
	}

	// Simpan sementara ke disk (opsional, bisa juga pakai io.Reader langsung)
	tmpFile, err := os.CreateTemp("", "upload-*.xlsx")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create temp file"})
		return
	}
	defer os.Remove(tmpFile.Name()) // hapus file setelah selesai

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
	tmpFile.Close() // harus di-close sebelum open oleh excelize

	// Baca Excel pakai excelize
	f, err := excelize.OpenFile(tmpFile.Name())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Excel file"})
		return
	}
	defer f.Close()

	// Ambil data dari sheet pertama
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not read rows"})
		return
	}

	// Contoh: Tampilkan isi baris sebagai array of array string
	c.JSON(http.StatusOK, gin.H{
		"sheet": sheetName,
		"rows":  rows,
	})
}

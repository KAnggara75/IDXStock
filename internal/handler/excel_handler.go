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
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/handler
 */

package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func StockList(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File not found"})
	}
	filename := fileHeader.Filename

	// Validasi ekstensi .xlsx
	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File harus .xlsx"})
	}

	// Validasi prefix (bisa configurable pakai viper)
	const prefix = "Daftar Saham"
	if !strings.HasPrefix(filename, prefix) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama file harus diawali dengan 'Daftar Saham'"})
	}

	// Validasi tanggal di belakang (YYYYMMDD)
	datePattern := regexp.MustCompile(`(\d{8})\.xlsx$`)
	matches := datePattern.FindStringSubmatch(filename)
	if matches == nil || len(matches) < 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Nama file harus diakhiri dengan tanggal format YYYYMMDD sebelum .xlsx"})
	}
	tanggalFile := matches[1]
	today := time.Now().Format("20060102")
	if tanggalFile != today {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":      fmt.Sprintf("Tanggal di filename harus hari ini: %s", today),
			"file_date":  tanggalFile,
			"today_date": today,
		})
	}

	// Simpan file ke tmp
	tempPath := filepath.Join(os.TempDir(), filename)
	if err := c.SaveFile(fileHeader, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not save temp file"})
	}
	defer func() {
		if err := os.Remove(tempPath); err != nil {
			fmt.Printf("Warning: failed to remove %s: %v\n", tempPath, err)
		}
	}()

	// Baca isi Excel
	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Excel file"})
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("Warning: failed to close excel file: %v\n", err)
		}
	}()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not read rows"})
	}

	return c.JSON(fiber.Map{
		"filename":  filename,
		"sheet":     sheetName,
		"rows":      rows,
		"validated": true,
	})
}

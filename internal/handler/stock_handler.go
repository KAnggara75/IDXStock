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
 * @author KAnggara75 on Wed 25/06/25 07.57
 * @project IDXStock handler
 * https://github.com/KAnggara75/IDXStock/tree/main/handler
 */

package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Stock struct {
	Code          string `json:"code"`
	Name          string `json:"name"`
	ListingDate   string `json:"listing_date,omitempty"`
	DelistingDate string `json:"delisting_date,omitempty"`
	Shares        int64  `json:"shares"`
	Board         string `json:"board"`
}

func UploadStocks(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file not found"})
	}
	filename := fileHeader.Filename

	if !strings.HasSuffix(strings.ToLower(filename), ".xlsx") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "file harus .xlsx"})
	}

	tempPath := filepath.Join(os.TempDir(), filename)
	if err := c.SaveFile(fileHeader, tempPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not save temp file"})
	}
	defer os.Remove(tempPath)

	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid Excel file"})
	}
	defer func() {
		_ = f.Close()
	}()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not read rows or sheet kosong"})
	}

	// Asumsi header kolom urut: Code, Name, ListingDate, DelistingDate, Shares, Board
	var stocks []Stock
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		// minimal ada 6 kolom, sesuaikan kalau excel-mu beda
		if len(row) < 6 {
			continue
		}
		shares := int64(0)
		if s := strings.TrimSpace(row[4]); s != "" {
			_, _ = fmt.Sscanf(s, "%d", &shares)
		}

		stocks = append(stocks, Stock{
			Code:          strings.TrimSpace(row[0]),
			Name:          strings.TrimSpace(row[1]),
			ListingDate:   parseDate(row[2]),
			DelistingDate: parseDate(row[3]),
			Shares:        shares,
			Board:         strings.TrimSpace(row[5]),
		})
	}

	return c.JSON(stocks)
}

// Helper parse date, biar format konsisten (YYYY-MM-DD), atau kosong kalau gagal
func parseDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// Coba format YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("2006-01-02")
	}
	return s // fallback: tetap return original string
}

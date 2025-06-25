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
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

type Stock struct {
	Code         string `json:"code"`
	CompanyName  string `json:"company_name"`
	ListingDate  string `json:"listing_date"`
	Shares       int64  `json:"shares"`
	ListingBoard string `json:"listing_board"`
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

	// Cari header index
	header := rows[0]
	idxCode := findIndex(header, "Code")
	idxCompany := findIndex(header, "Company Name")
	idxListingDate := findIndex(header, "Listing Date")
	idxShares := findIndex(header, "Shares")
	idxBoard := findIndex(header, "Listing Board")

	var stocks []Stock
	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}
		if idxCode == -1 || idxCompany == -1 || idxListingDate == -1 || idxShares == -1 || idxBoard == -1 {
			continue
		}
		if len(row) <= idxBoard {
			continue
		}

		shares := parseShares(getOrEmpty(row, idxShares))

		stocks = append(stocks, Stock{
			Code:         getOrEmpty(row, idxCode),
			CompanyName:  getOrEmpty(row, idxCompany),
			ListingDate:  parseDate(getOrEmpty(row, idxListingDate)),
			Shares:       shares,
			ListingBoard: getOrEmpty(row, idxBoard),
		})
	}

	return c.JSON(stocks)
}

func findIndex(header []string, name string) int {
	for i, h := range header {
		if strings.EqualFold(strings.TrimSpace(h), strings.TrimSpace(name)) {
			return i
		}
	}
	return -1
}

func getOrEmpty(row []string, idx int) string {
	if idx >= 0 && idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

func parseDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("2006-01-02")
	}
	return s
}

func parseShares(s string) int64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		// Kalau error, bisa juga coba parse float, lalu dibulatkan ke int64
		if f, err2 := strconv.ParseFloat(s, 64); err2 == nil {
			return int64(f)
		}
		return 0
	}
	return val
}

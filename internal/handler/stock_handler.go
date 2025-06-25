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
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/handler
 */

package handler

import (
	"fmt"
	"github.com/KAnggara75/IDXStock/internal/excel"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strings"
)

type Stock struct {
	Code         string `json:"code"`
	CompanyName  string `json:"company_name"`
	ListingDate  string `json:"listing_date"`
	Shares       int64  `json:"shares"`
	ListingBoard string `json:"listing_board"`
}

func ConvertStocks(c *fiber.Ctx) error {
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
	defer func() {
		if err := os.Remove(tempPath); err != nil {
			fmt.Printf("Warning: failed to remove %s: %v\n", tempPath, err)
		}
	}()

	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid Excel file"})
	}
	defer func() { _ = f.Close() }()

	stocks, err := excel.ParseStock(f)

	return c.JSON(stocks)
}

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
	"github.com/KAnggara75/IDXStock/internal/helper"
	"github.com/KAnggara75/IDXStock/internal/utils"
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

var (
	// ENGLISH HEADERS
	headerEng = map[string]string{
		"code":          "Code",
		"company_name":  "Company Name",
		"listing_date":  "Listing Date",
		"shares":        "Shares",
		"listing_board": "Listing Board",
	}
	// INDONESIAN HEADERS
	headerIndo = map[string]string{
		"code":          "Kode",
		"company_name":  "Nama Perusahaan",
		"listing_date":  "Tanggal Pencatatan",
		"shares":        "Saham",
		"listing_board": "Papan Pencatatan",
	}
)

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

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not read rows or sheet kosong"})
	}

	header := rows[0]
	var headerMap map[string]string
	if utils.FindIndex(header, headerEng["code"]) != -1 && utils.FindIndex(header, headerEng["company_name"]) != -1 {
		headerMap = headerEng
	} else if utils.FindIndex(header, headerIndo["code"]) != -1 && utils.FindIndex(header, headerIndo["company_name"]) != -1 {
		headerMap = headerIndo
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Header tidak dikenali (harus Inggris atau Indonesia)"})
	}

	idxCode := utils.FindIndex(header, headerMap["code"])
	idxShares := utils.FindIndex(header, headerMap["shares"])
	idxBoard := utils.FindIndex(header, headerMap["listing_board"])
	idxCompany := utils.FindIndex(header, headerMap["company_name"])
	idxListingDate := utils.FindIndex(header, headerMap["listing_date"])

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

		stocks = append(stocks, Stock{
			Code:         utils.GetOrEmpty(row, idxCode),
			CompanyName:  utils.GetOrEmpty(row, idxCompany),
			ListingDate:  utils.ParseDateFlexible(utils.GetOrEmpty(row, idxListingDate)),
			Shares:       utils.ParseToNumber(utils.GetOrEmpty(row, idxShares)),
			ListingBoard: helper.MapBoardToEN(utils.GetOrEmpty(row, idxBoard)),
		})
	}

	return c.JSON(stocks)
}

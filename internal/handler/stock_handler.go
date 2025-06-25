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
	boardIDtoEN = map[string]string{
		"Akselerasi":        "Acceleration",
		"Pengembangan":      "Development",
		"Ekonomi Baru":      "Ekonomi Baru",
		"Utama":             "Main",
		"Pemantauan Khusus": "Watchlist",
		"Acceleration":      "Acceleration",
		"Development":       "Development",
		"Main":              "Main",
		"Watchlist":         "Watchlist",
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
	defer os.Remove(tempPath)

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
	if findIndex(header, headerEng["code"]) != -1 && findIndex(header, headerEng["company_name"]) != -1 {
		headerMap = headerEng
	} else if findIndex(header, headerIndo["code"]) != -1 && findIndex(header, headerIndo["company_name"]) != -1 {
		headerMap = headerIndo
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Header tidak dikenali (harus Inggris atau Indonesia)"})
	}

	idxCode := findIndex(header, headerMap["code"])
	idxShares := findIndex(header, headerMap["shares"])
	idxBoard := findIndex(header, headerMap["listing_board"])
	idxCompany := findIndex(header, headerMap["company_name"])
	idxListingDate := findIndex(header, headerMap["listing_date"])

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

		shares := parseSharesMulti(getOrEmpty(row, idxShares))

		stocks = append(stocks, Stock{
			Code:         getOrEmpty(row, idxCode),
			CompanyName:  getOrEmpty(row, idxCompany),
			ListingDate:  parseDateFlexible(getOrEmpty(row, idxListingDate)),
			Shares:       shares,
			ListingBoard: mapBoardToEN(getOrEmpty(row, idxBoard)),
		})
	}

	return c.JSON(stocks)
}

// --- UTILS ---

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

// Parse format tanggal Indonesia & Inggris (contoh: "09 Dec 1997", "09 Des 1997", fallback as is)
func parseDateFlexible(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	formats := []string{
		"02 Jan 2006", // English
		"02 Jan 06",
		"02 January 2006",
		"02 Jan 2006", // Bisa juga dipakai untuk "09 Dec 1997"
		"02 Jan 2006", // Repeated so that code is easy to extend
	}

	idToEnMonth := map[string]string{
		"Jan": "Jan", "Feb": "Feb", "Mar": "Mar", "Apr": "Apr", "Mei": "May",
		"Jun": "Jun", "Jul": "Jul", "Agu": "Aug", "Sep": "Sep", "Okt": "Oct",
		"Nov": "Nov", "Des": "Dec",
	}

	parts := strings.Split(s, " ")
	if len(parts) == 3 {
		if m, ok := idToEnMonth[parts[1]]; ok {
			parts[1] = m
			s = strings.Join(parts, " ")
		}
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return s
}

func parseSharesMulti(s string) int64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.ReplaceAll(s, ".", "")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	val, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if f, err2 := strconv.ParseFloat(s, 64); err2 == nil {
			return int64(f)
		}
		return 0
	}
	return val
}

func mapBoardToEN(s string) string {
	s = strings.TrimSpace(s)
	if en, ok := boardIDtoEN[s]; ok {
		return en
	}
	// Kalau tidak cocok, return as-is
	return s
}

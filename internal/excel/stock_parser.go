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
 * @author KAnggara75 on Wed 25/06/25 21.00
 * @project IDXStock excel
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/excel
 */

package excel

import (
	"errors"
	"github.com/KAnggara75/IDXStock/internal/domain"
	"github.com/KAnggara75/IDXStock/internal/helper"
	"github.com/KAnggara75/IDXStock/internal/utils"
	"github.com/xuri/excelize/v2"
	"log"
)

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

func ParseStock(f *excelize.File) ([]domain.Stock, error) {
	log.Printf("[WARN] ParseStock")
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil || len(rows) < 2 {
		return nil, errors.New("could not read rows or empty sheet")
	}

	header := rows[0]
	var headerMap map[string]string
	switch {
	case utils.FindIndex(header, headerEng["code"]) != -1 && utils.FindIndex(header, headerEng["company_name"]) != -1:
		headerMap = headerEng
	case utils.FindIndex(header, headerIndo["code"]) != -1 && utils.FindIndex(header, headerIndo["company_name"]) != -1:
		headerMap = headerIndo
	default:
		return nil, errors.New("header tidak found")
	}

	idxCode := utils.FindIndex(header, headerMap["code"])
	idxShares := utils.FindIndex(header, headerMap["shares"])
	idxBoard := utils.FindIndex(header, headerMap["listing_board"])
	idxCompany := utils.FindIndex(header, headerMap["company_name"])
	idxListingDate := utils.FindIndex(header, headerMap["listing_date"])

	if idxCode == -1 || idxCompany == -1 || idxListingDate == -1 || idxShares == -1 || idxBoard == -1 {
		return nil, errors.New("header not enough or not found")
	}

	var stocks []domain.Stock
	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) <= idxBoard {
			continue
		}

		stocks = append(stocks, domain.Stock{
			Code:         utils.GetOrEmpty(row, idxCode),
			CompanyName:  utils.GetOrEmpty(row, idxCompany),
			ListingDate:  utils.ParseDateFlexible(utils.GetOrEmpty(row, idxListingDate)),
			Shares:       utils.ParseToNumber(utils.GetOrEmpty(row, idxShares)),
			ListingBoard: helper.MapBoardToEN(utils.GetOrEmpty(row, idxBoard)),
		})
	}

	return stocks, nil
}

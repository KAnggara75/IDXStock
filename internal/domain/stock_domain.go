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
 * @author KAnggara75 on Thu 26/06/25 22.42
 * @project IDXStock domain
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/domain
 */

package domain

import "github.com/KAnggara75/IDXStock/internal/dto"

type Stock struct {
	Code          string
	CompanyName   string
	ListingDate   string
	DelistingDate string
	ListingBoard  string
	Shares        int64
}

func ToDTO(stock Stock) dto.Stock {
	return dto.Stock{
		Code:          stock.Code,
		CompanyName:   stock.CompanyName,
		ListingDate:   stock.ListingDate,
		DelistingDate: stock.DelistingDate,
		ListingBoard:  stock.ListingBoard,
		Shares:        stock.Shares,
	}
}

func SliceToDTO(stocks []Stock) []dto.Stock {
	result := make([]dto.Stock, len(stocks))
	for i, s := range stocks {
		result[i] = ToDTO(s)
	}
	return result
}

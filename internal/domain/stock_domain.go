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

import (
	"github.com/KAnggara75/IDXStock/internal/dto"
	"github.com/KAnggara75/IDXStock/internal/logx"
	"github.com/KAnggara75/IDXStock/internal/repository/model"
	"time"
)

type Stock struct {
	Code          string
	CompanyName   string
	ListingDate   string
	DelistingDate string
	ListingBoard  string
	Shares        int64
}

func ToStockDTO(stock Stock) dto.Stock {
	return dto.Stock{
		Code:          stock.Code,
		CompanyName:   stock.CompanyName,
		ListingDate:   stock.ListingDate,
		DelistingDate: stock.DelistingDate,
		ListingBoard:  stock.ListingBoard,
		Shares:        stock.Shares,
	}
}

func ToStockDTOs(stocks []Stock) []dto.Stock {
	logx.Debug("Invoke ToStockDTOs function for Stock domain")
	result := make([]dto.Stock, len(stocks))
	for i, s := range stocks {
		result[i] = ToStockDTO(s)
	}
	return result
}

func ToStockModel(d Stock) (model.Stock, error) {
	const layout = "2006-01-02"
	listingDate, err := time.Parse(layout, d.ListingDate)
	if err != nil {
		logx.Errorf("Failed to parse listing date: %v", err)
		return model.Stock{}, err
	}
	var delistingDate *time.Time
	if d.DelistingDate != "" {
		dd, err := time.Parse(layout, d.DelistingDate)
		if err != nil {
			logx.Errorf("Failed to parse delisting date: %v", err)
			return model.Stock{}, err
		}
		delistingDate = &dd
	}
	return model.Stock{
		Code:          d.Code,
		Name:          d.CompanyName,
		ListingDate:   listingDate,
		DelistingDate: delistingDate,
		Shares:        d.Shares,
		Board:         d.ListingBoard,
	}, nil
}

func ToStockModels(domains []Stock) ([]*model.Stock, []error) {
	result := make([]*model.Stock, 0, len(domains))
	var errs []error
	for _, d := range domains {
		m, err := ToStockModel(d)
		if err != nil {
			logx.Errorf("Failed to convert Stock domain to model, SKIPPING: %v", err)
			errs = append(errs, err)
			continue
		}
		result = append(result, &m)
	}
	return result, errs
}

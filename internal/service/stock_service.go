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
 * @author KAnggara75 on Sun 22/06/25 21.29
 * @project IDXStock service
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/service
 */

package service

import (
	"github.com/KAnggara75/IDXStock/internal/repository"
	"github.com/KAnggara75/IDXStock/internal/repository/model"
)

type StockService struct {
	repo repository.StockRepository
}

func NewStockService(r repository.StockRepository) *StockService {
	return &StockService{repo: r}
}

func (s *StockService) UpsertStocks(stocks []*model.Stock) error {
	return s.repo.UpsertStocks(stocks)
}

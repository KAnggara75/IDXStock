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
 * @author KAnggara75 on Thu 26/06/25 05.48
 * @project IDXStock repository
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/repository
 */

package repository

import (
	"github.com/KAnggara75/IDXStock/internal/logx"
	"github.com/KAnggara75/IDXStock/internal/repository/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) UpsertStock(stock *model.Stock) error {
	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "code"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"name", "listing_date", "delisting_date", "shares", "board", "last_modified",
		}),
	}).Create(&stock).Error
}

func (r *stockRepository) UpsertStocks(stocks []*model.Stock) error {
	for _, stock := range stocks {
		if err := r.UpsertStock(stock); err != nil {
			logx.Errorf("Failed to upsert stocks: %v", err)
			return err
		}
	}
	return nil
}

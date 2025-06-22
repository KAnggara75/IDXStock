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
 * @author KAnggara75 on Sun 22/06/25 21.30
 * @project IDXStock model
 * https://github.com/KAnggara75/IDXStock/tree/main/model
 */

package model

import (
	"time"
)

type Stock struct {
	Code          string     `gorm:"primaryKey;column:code;type:varchar(10)"`
	Name          string     `gorm:"column:name;type:varchar(200);not null"`
	ListingDate   *time.Time `gorm:"column:listing_date"`
	DelistingDate *time.Time `gorm:"column:delisting_date"`
	Shares        int64      `gorm:"column:shares;not null"`
	Board         string     `gorm:"column:board;type:idxstock.board;default:Utama;not null"`
	LastModified  time.Time  `gorm:"column:last_modified;autoUpdateTime"`
}

func (Stock) TableName() string {
	return "idxstock.stocks"
}

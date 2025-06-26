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
 * @author KAnggara75 on Thu 26/06/25 23.50
 * @project IDXStock app
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/app
 */

package app

import (
	"github.com/KAnggara75/IDXStock/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConn() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.GetDBConn()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(config.GetDBMaxOpen())
	sqlDB.SetMaxIdleConns(config.GetDBMaxIdle())
	sqlDB.SetConnMaxLifetime(config.GetDBMaxLifetime())
	return db, nil
}

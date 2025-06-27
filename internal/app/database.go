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
	"fmt"
	"github.com/KAnggara75/IDXStock/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDBConn() (*gorm.DB, error) {
	dsn := config.GetDBConn()
	NewLogger().Debug().Str("trace_id", config.Get40Space()).Msgf("Connecting to database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	NewLogger().Debug().Str("trace_id", config.Get40Space()).Msgf("DB connection established successfully.")

	// Set connection pool settings
	NewLogger().Debug().Str("trace_id", config.Get40Space()).Msgf("Setting connection pool parameters")
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB from GORM: %w", err)
	}
	NewLogger().Debug().Str("trace_id", config.Get40Space()).Msgf("SQL DB connection pool obtained successfully.")

	maxOpen := config.GetDBMaxOpen()
	maxIdle := config.GetDBMaxIdle()
	maxLifetime := config.GetDBMaxLifetime()

	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(maxLifetime)

	NewLogger().Debug().Str("trace_id", config.Get40Space()).Msgf("Pinging database to ensure connection is alive")
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	NewLogger().Debug().Str("trace_id", config.Get40Space()).Msgf("Database connection is alive and ready to use.")

	return db, nil
}

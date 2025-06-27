/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/>.
 *
 * @author KAnggara75 on Sat 21/06/25 23.59
 * @project IDXStock main
 * https://github.com/KAnggara75/IDXStock/tree/main/cmd
 */

package main

import (
	"github.com/KAnggara75/IDXStock/internal/app"
	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/KAnggara75/IDXStock/internal/logx"
	"github.com/KAnggara75/IDXStock/internal/route"
	"github.com/KAnggara75/IDXStock/internal/utils"
	"github.com/KAnggara75/scc2go"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	scc2go.GetEnv(os.Getenv("SCC_IDXSTOCK_URL"), os.Getenv("AUTH"))
}

func main() {
	logx.Info("Starting IDXStock application...")
	db, err := app.NewDBConn()
	if err != nil {
		logx.Fatalf("failed to connect DB: %v", err)
	} else {
		logx.Debug("DB connection established successfully.")
	}

	sqlDB, err := db.DB()
	if err != nil {
		logx.Fatalf("failed to get *sql.DB: %v", err)
	}
	defer func() {
		logx.Info("Closing DB connection...")
		if err := sqlDB.Close(); err != nil {
			logx.Errorf("error closing DB: %v", err)
		}
		logx.Info("DB connection closed.")
	}()

	if config.IsAutoMigrate() {
		logx.Info("DB AutoMigrate is enabled. Running migration...")
		if err := utils.AutoMigrate(db); err != nil {
			logx.Fatalf("Migration failed: %v", err)
		}
	} else {
		logx.Info("DB AutoMigrate is disabled. Skipping migration.")
	}

	server := route.NewRouter()
	go func() {
		if err := server.Listen(config.GetPort()); err != nil {
			logx.Fatalf("server error: %v", err)
		}
	}()

	logx.Infof("Server started at %s", config.GetPort())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logx.Warn("Shutting down server...")
	logx.Warn("App stopped gracefully.")
}

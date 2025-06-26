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
	"github.com/KAnggara75/IDXStock/internal/route"
	"github.com/KAnggara75/IDXStock/internal/utils"
	"github.com/KAnggara75/scc2go"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	scc2go.GetEnv(os.Getenv("SCC_IDXSTOCK_URL"), os.Getenv("AUTH"))
}

func main() {
	db, err := app.NewDBConn()
	if err != nil {
		panic(err)
	}

	if config.IsAutoMigrate() {
		log.Println("DB AutoMigrate is enabled. Running migration...")
		if err := utils.AutoMigrate(db); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Println("DB AutoMigrate is disabled. Skipping migration.")
	}

	server := route.NewRouter()
	go func() {
		if err := server.Listen(config.GetPort()); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Printf("Server started at %s", config.GetPort())

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	log.Println("App stopped gracefully.")
}

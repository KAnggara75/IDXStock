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
 * @author KAnggara75 on Sun 22/06/25 15.44
 * @project IDXStock route
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/route
 */

package route

import (
	"github.com/KAnggara75/IDXStock/internal/handler"
	"github.com/KAnggara75/IDXStock/internal/logx"
	"github.com/gofiber/fiber/v2"
)

type HandlerSet struct {
	Stock *handler.StockHandler
}

func NewRouter(h *HandlerSet) *fiber.App {
	logx.Debug("Initializing Fiber router...")
	fiberApp := fiber.New()
	fiberApp.Post("/upload", h.Stock.ConvertStocks)
	return fiberApp
}

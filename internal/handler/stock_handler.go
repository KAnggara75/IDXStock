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
 * @author KAnggara75 on Wed 25/06/25 07.57
 * @project IDXStock handler
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/handler
 */

package handler

import (
	"github.com/KAnggara75/IDXStock/internal/domain"
	"github.com/KAnggara75/IDXStock/internal/excel"
	"github.com/KAnggara75/IDXStock/internal/helper"
	"github.com/KAnggara75/IDXStock/internal/logx"
	"github.com/KAnggara75/IDXStock/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
	"os"
)

type StockHandler struct {
	service *service.StockService
}

func NewStockHandler(s *service.StockService) *StockHandler {
	return &StockHandler{service: s}
}

func (h *StockHandler) ConvertStocks(c *fiber.Ctx) error {
	logx.Debug("Invoke ConvertStocks handler")
	tempPath, err := helper.SaveTempFile(c)
	if err != nil {
		logx.Errorf("SaveTempFile Failed: %v", err)
		return helper.FiberErr(c, fiber.StatusBadRequest, err.Error())
	}
	defer func() { _ = os.Remove(tempPath) }()

	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		logx.Errorf("OpenFile Failed %s: %v", tempPath, err)
		return helper.FiberErr(c, fiber.StatusBadRequest, "invalid Excel file")
	}
	defer func() { _ = f.Close() }()

	stocks, err := excel.ParseStock(f)
	if err != nil {
		logx.Errorf("ParseStock Failed: %v", err)
		return helper.FiberErr(c, fiber.StatusBadRequest, err.Error())
	}

	go func(data []domain.Stock) {
		stockModels, errs := domain.ToStockModels(data)
		if len(errs) > 0 {
			for _, err := range errs {
				logx.Warnf("Skipping invalid record: %v", err)
			}
		}
		if upsertErr := h.service.UpsertStocks(stockModels); upsertErr != nil {
			logx.Warnf("Failed to upsert stocks: %v", upsertErr)
		}
	}(stocks)

	response := domain.ToStockDTOs(stocks)
	return c.JSON(response)
}

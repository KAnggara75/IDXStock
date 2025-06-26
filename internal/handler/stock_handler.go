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
	"github.com/KAnggara75/IDXStock/internal/helper"
	"os"
	"path/filepath"
	"strings"

	"github.com/KAnggara75/IDXStock/internal/excel"
	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

func ConvertStocks(c *fiber.Ctx) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return helper.FiberErr(c, fiber.StatusBadRequest, "file not found")
	}

	if ext := strings.ToLower(filepath.Ext(fileHeader.Filename)); ext != ".xlsx" {
		return helper.FiberErr(c, fiber.StatusBadRequest, "file harus .xlsx")
	}

	tempPath := filepath.Join(os.TempDir(), fileHeader.Filename)
	if err := c.SaveFile(fileHeader, tempPath); err != nil {
		return helper.FiberErr(c, fiber.StatusInternalServerError, "could not save temp file")
	}
	defer func() { _ = os.Remove(tempPath) }()

	f, err := excelize.OpenFile(tempPath)
	if err != nil {
		return helper.FiberErr(c, fiber.StatusBadRequest, "invalid Excel file")
	}
	defer func() { _ = f.Close() }()

	stocks, err := excel.ParseStock(f)
	if err != nil {
		return helper.FiberErr(c, fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(stocks)
}

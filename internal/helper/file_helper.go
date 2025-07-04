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
 * @author KAnggara75 on Fri 27/06/25 00.58
 * @project IDXStock helper
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/helper
 */

package helper

import (
	"github.com/KAnggara75/IDXStock/internal/logx"
	"github.com/KAnggara75/IDXStock/internal/validation"
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
)

func SaveTempFile(c *fiber.Ctx) (string, error) {
	logx.Debug("Invoke SaveTempFile helper")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		logx.Errorf("Failed to get form file: %v", err)
		return "", err
	}

	if err := validation.ValidExt(fileHeader.Filename, ".xlsx"); err != nil {
		logx.Errorf("Invalid file extension: %v", err)
		return "", err
	}

	tempPath := filepath.Join(os.TempDir(), fileHeader.Filename)
	if err := c.SaveFile(fileHeader, tempPath); err != nil {
		logx.Errorf("Failed to save file to temp path %s: %v", tempPath, err)
		return "", err
	}
	return tempPath, nil
}

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
 * @author KAnggara75 on Fri 27/06/25 00.47
 * @project IDXStock validation
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/validation
 */

package validation

import (
	"errors"
	"path/filepath"
	"strings"
)

func ValidExt(filename string, allowedExtensions ...string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, a := range allowedExtensions {
		if ext == strings.ToLower(a) {
			return nil
		}
	}
	return errors.New("file must: " + strings.Join(allowedExtensions, ", "))
}

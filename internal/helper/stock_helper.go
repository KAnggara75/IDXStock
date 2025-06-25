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
 * @author KAnggara75 on Wed 25/06/25 20.55
 * @project IDXStock helper
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/helper
 */

package helper

import "strings"

var boardIDtoEN = map[string]string{
	"Main":              "Main",
	"Utama":             "Main",
	"Watchlist":         "Watchlist",
	"Akselerasi":        "Acceleration",
	"Development":       "Development",
	"Pengembangan":      "Development",
	"Ekonomi Baru":      "Ekonomi Baru",
	"Acceleration":      "Acceleration",
	"Pemantauan Khusus": "Watchlist",
}

func MapBoardToEN(s string) string {
	s = strings.TrimSpace(s)
	if en, ok := boardIDtoEN[s]; ok {
		return en
	}
	return s
}

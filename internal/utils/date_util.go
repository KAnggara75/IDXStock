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
 * @author KAnggara75 on Wed 25/06/25 09.22
 * @project IDXStock utils
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/utils
 */

package utils

import (
	"strings"
	"time"
)

func ParseDateFlexible(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	formats := []string{
		"02 Jan 06",
		"02 Jan 2006",
		"02 January 2006",
	}

	idToEnMonth := map[string]string{
		"Jan": "Jan", "Feb": "Feb", "Mar": "Mar", "Apr": "Apr", "Mei": "May",
		"Jun": "Jun", "Jul": "Jul", "Agu": "Aug", "Sep": "Sep", "Okt": "Oct",
		"Nov": "Nov", "Des": "Dec",
	}

	parts := strings.Split(s, " ")
	if len(parts) == 3 {
		if m, ok := idToEnMonth[parts[1]]; ok {
			parts[1] = m
			s = strings.Join(parts, " ")
		}
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return s
}

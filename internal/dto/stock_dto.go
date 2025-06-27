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
 * @author KAnggara75 on Wed 25/06/25 21.40
 * @project IDXStock dto
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/dto
 */

package dto

type Stock struct {
	Code          string `json:"code"`
	CompanyName   string `json:"company_name"`
	ListingDate   string `json:"listing_date"`
	DelistingDate string `json:"delisting_date"`
	Shares        int64  `json:"shares"`
	ListingBoard  string `json:"listing_board"`
}

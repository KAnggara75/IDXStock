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
 * @author KAnggara75 on Wed 25/06/25 21.46
 * @project IDXStock repository
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/repository
 */

package repository

import (
	"github.com/KAnggara75/IDXStock/internal/repository/model"
)

type StockRepository interface {
	UpsertStock(stock *model.Stock) error
}

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
 * @author KAnggara75 on Sun 22/06/25 21.29
 * @project IDXStock service
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/service
 */

package service

import (
	"github.com/KAnggara75/IDXStock/internal/domain"
	"github.com/KAnggara75/IDXStock/internal/logx"
)

func UpsertStocks(stocks []domain.Stock) error {
	logx.Debug("Invoke UpsertStocks service")
	//for _, s := range stocks {
	//if err := repository.UpsertStock(s); err != nil {
	//	return err
	//}
	//}
	return nil
}

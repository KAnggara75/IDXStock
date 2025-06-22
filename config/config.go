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
 * @author KAnggara75 on Sun 22/06/25 00.39
 * @project IDXStock config
 * https://github.com/KAnggara75/IDXStock/tree/main/config
 */

package config

import (
	"github.com/spf13/viper"
)

func GetDBConn() string {
	return viper.GetString("db.idxstock.host")
}

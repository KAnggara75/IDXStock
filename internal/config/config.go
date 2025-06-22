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
	"strings"
)

func GetDBConn() string {
	return viper.GetString("db.idxstock.host")
}

func GetStockListPrefix() string {
	prefix := viper.GetString("app.stockListPrefix")
	if prefix == "" {
		prefix = "Daftar Saham"
	}
	return prefix
}

func GetGinMode() string {
	mode := viper.GetString("gin.mode")
	switch mode {
	case "debug", "":
		return "debug"
	case "release":
		return "release"
	case "test":
		return "test"
	default:
		return "debug"
	}
}

func GetTrustedProxies() []string {
	proxies := viper.GetString("gin.trusted_proxies")
	if proxies == "" {
		return []string{"127.0.0.1"}
	}
	return strings.Split(proxies, ",")
}

func GetPort() string {
	port := viper.GetString("gin.port")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return port
}

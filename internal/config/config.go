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
 * https://github.com/KAnggara75/IDXStock/tree/main/internal/config
 */

package config

import (
	"github.com/spf13/viper"
	"strings"
	"time"
)

const (
	defaultDBMaxIdle     = 1
	defaultDBMaxOpen     = 10
	defaultDBMaxLifetime = 10 * time.Minute
)

func GetDBConn() string {
	return viper.GetString("db.idxstock.host")
}

func GetDBMaxIdle() int {
	return getIntConfig("db.idxstock.maxIdle", defaultDBMaxIdle)
}

func GetDBMaxOpen() int {
	return getIntConfig("db.idxstock.maxOpen", defaultDBMaxOpen)
}

func GetDBMaxLifetime() time.Duration {
	return getDurationConfig("db.idxstock.maxLifetime", defaultDBMaxLifetime)
}

func IsAutoMigrate() bool {
	return viper.GetBool("db.autoMigrate")
}

func GetLogLevel() string { return viper.GetString("log.level") }

func Get40Space() string {
	return strings.Repeat(" ", 40)
}

func GetPort() string {
	port := viper.GetString("app.port")
	if port == "" {
		port = "8080"
	}
	if port[0] != ':' {
		port = ":" + port
	}
	return port
}

func getIntConfig(key string, def int) int {
	if viper.IsSet(key) {
		val := viper.GetInt(key)
		if val > 0 {
			return val
		}
	}
	return def
}

func getDurationConfig(key string, def time.Duration) time.Duration {
	if viper.IsSet(key) {
		val := viper.GetDuration(key)
		if val > 0 {
			return val
		}
	}
	return def
}

/*
 * Copyright (c) 2025 KAnggara75
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * See <https://www.gnu.org/licenses/>.
 *
 * @author KAnggara75 on Sat 21/06/25 23.59
 * @project IDXStock main
 * https://github.com/KAnggara75/IDXStock/tree/main/main
 */

package main

import (
	"fmt"
	"github.com/KAnggara75/IDXStock/config"
	"github.com/KAnggara75/scc2go"
	"github.com/gin-gonic/gin"
	"os"
)

func init() {
	scc2go.GetEnv(os.Getenv("SCC_IDXSTOCK_URL"), os.Getenv("AUTH"))
}

func main() {
	gin.SetMode(config.GetGinMode())

	router := gin.Default()
	if err := router.SetTrustedProxies(config.GetTrustedProxies()); err != nil {
		panic(err)
	}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	if err := router.Run(config.GetPort()); err != nil {
		panic(err)
	}

	fmt.Println("IDXStock server is running on port", config.GetPort())
}

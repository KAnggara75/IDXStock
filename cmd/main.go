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
	scc2go.GetEnv(os.Getenv("SCC_URL"), os.Getenv("AUTH"))
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		deviceId := config.GetDBConn()
		fmt.Println("Device ID:", deviceId)
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := router.Run()
	if err != nil {
		return
	} // listen and serve on 0.0.0.0:8080
}

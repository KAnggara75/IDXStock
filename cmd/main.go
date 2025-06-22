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
	"github.com/KAnggara75/IDXStock/internal/config"
	"github.com/KAnggara75/IDXStock/internal/route"
	"github.com/KAnggara75/scc2go"
	"os"
)

func init() {
	scc2go.GetEnv(os.Getenv("SCC_IDXSTOCK_URL"), os.Getenv("AUTH"))
}

func main() {
	app := route.SetupRouter()
	if err := app.Listen(config.GetPort()); err != nil {
		panic(err)
	}
}

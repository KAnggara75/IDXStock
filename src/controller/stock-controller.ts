import { Hono } from "hono";
import type { User } from "@prisma/client";
import type { ApplicationVariables } from "../model/app-model";
import { StockUsecase } from "../usecase/stock-usecase.ts";

export const stockController = new Hono<{ Variables: ApplicationVariables }>();

stockController.patch("/stocks/google", async (c) => {
	const user: User = c.get("user");
	const response = await StockUsecase.updateDailyStock(user);

	return c.json({
		data: response,
	});
});

stockController.post("/stocks/idx", async (c) => {
	const user: User = c.get("user");

	const body = await c.req.parseBody();

	const xlxsFile: string | File = body.file;

	// Convert xlsx to JSON
	const response = await StockUsecase.addStockSummary(user, xlxsFile);

	return c.json({
		data: response,
	});
});

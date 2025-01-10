import { Hono } from "hono";
import type { ApplicationVariables } from "../model/app-model";
import { StockUpdateUsecase } from "../usecase/stockUpdate-usecase";

export const stockController = new Hono<{ Variables: ApplicationVariables }>();

stockController.get("/stock", async (c) => {
	const response = await StockUpdateUsecase.updateStockData();

	return c.json({
		data: response,
	});
});

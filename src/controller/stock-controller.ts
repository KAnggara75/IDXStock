import { Hono } from "hono";
import { log } from "../config/logger";
import type { User } from "@prisma/client";
import type { ApplicationVariables } from "../model/app-model";
import { authMiddleware } from "../middleware/auth-middleware";
import { StockUpdateUsecase } from "../usecase/stockUpdate-usecase";

export const stockController = new Hono<{ Variables: ApplicationVariables }>();

stockController.use(authMiddleware);

stockController.patch("/stocks", async (c) => {
	const user: User = c.get("user");
	log.info(JSON.stringify(user));
	const response = await StockUpdateUsecase.updateStockData(user);

	return c.json({
		data: response,
	});
});

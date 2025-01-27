import { Hono } from "hono";
import { log } from "../config/logger";
import type { User } from "@prisma/client";
import type { ApplicationVariables } from "../model/app-model";
import { StockUpdateUsecase } from "../usecase/stockUpdate-usecase";

export const stockController = new Hono<{ Variables: ApplicationVariables }>();

stockController.patch("/stocks/google", async (c) => {
	const user: User = c.get("user");
	log.info(JSON.stringify(user));
	const response = await StockUpdateUsecase.updateDailyStock(user);

	return c.json({
		data: response,
	});
});

stockController.post("/stocks/idx", async (c) => {
	const user: User = c.get("user");
	log.info(JSON.stringify(user));

	return c.json({
		data: "response",
	});
});

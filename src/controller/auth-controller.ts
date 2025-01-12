import { Hono } from "hono";
import { validator } from "hono/validator";
import { log } from "../config/logger.ts";
import { AuthService } from "../service/auth-service.ts";
import type { ApplicationVariables } from "../model/app-model.ts";
import { jsonMiddleware } from "../middleware/json-middleware.ts";

export const authController = new Hono<{ Variables: ApplicationVariables }>();

authController.use(jsonMiddleware);

authController.post(
	"/register",
	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const request = await c.req.json();

		const response = await AuthService.register(request);
		log.info(`Registering ${response.username}`);

		return c.json({
			data: response,
		});
	}
);

authController.post(
	"/login",

	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const request = await c.req.json();
		const response = await AuthService.login(request);

		return c.json({
			data: response,
		});
	}
);

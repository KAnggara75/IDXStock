import { Hono } from "hono";
import { validator } from "hono/validator";
import { log } from "../config/logger.ts";
import { UserService } from "../service/user-service.ts";
import type { ApplicationVariables } from "../model/app-model.ts";

export const authController = new Hono<{ Variables: ApplicationVariables }>();

authController.post(
	"/register",
	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const request = await c.req.json();

		const response = await UserService.register(request);
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
		const response = await UserService.login(request);

		return c.json({
			data: response,
		});
	}
);

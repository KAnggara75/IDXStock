import { Hono } from "hono";
import { log } from "../config/logger.ts";
import { validator } from "hono/validator";
import { AuthService } from "../service/auth-service.ts";
import type { ApplicationVariables } from "../model/app-model.ts";

export const authController = new Hono<{ Variables: ApplicationVariables }>();

authController.post(
	"/register",
	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const request = await c.req.json();

		const response = await AuthService.register(request);
		log.debug(`Registering ${response.username}`);

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

import { Hono } from "hono";
import { log } from "../config/logger.ts";
import type { User } from "@prisma/client";
import { validator } from "hono/validator";
import { AuthService } from "../service/auth-service.ts";
import type { ApplicationVariables } from "../model/app-model.ts";
import { authMiddleware } from "../middleware/auth-middleware.ts";

export const authController = new Hono<{ Variables: ApplicationVariables }>();

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

authController.use(authMiddleware);

authController.delete(
	"/logout",

	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const user: User = c.get("user");
		await AuthService.logout(user);
		return c.body(null, 204);
	}
);

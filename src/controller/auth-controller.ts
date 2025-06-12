import { Hono } from "hono";
import { log } from "../config/logger.ts";
import { AuthService } from "../service/auth-service.ts";
import type { UserResponse } from "../model/user-model.ts";
import type { ApplicationVariables } from "../model/app-model.ts";
import { AuthValidation } from "../validation/auth-validation.ts";
import { validateWithSchema } from "../validation/validate-with-schema.ts";

export const authController = new Hono<{ Variables: ApplicationVariables }>();

authController.post(
	"/register",
	validateWithSchema("json", AuthValidation.REGISTER),
	async (c) => {
		log.debug("Registering new user");
		const request = await c.req.json();
		const response: UserResponse = await AuthService.register(request);
		log.debug(`Registering ${response.username}`);

		return c.json({
			data: response,
		});
	}
);

authController.post(
	"/login",
	validateWithSchema("json", AuthValidation.LOGIN),
	async (c) => {
		const request = await c.req.json();
		const response = await AuthService.login(request);

		return c.json({
			data: response,
		});
	}
);

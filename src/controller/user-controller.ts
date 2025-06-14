import { Hono } from "hono";
import { log } from "../config/logger";
import { validator } from "hono/validator";
import type { User } from "@prisma/client";
import { UserService } from "../service/user-service";
import { AuthService } from "../service/auth-service";
import type { ApplicationVariables } from "../model/app-model";
import { type UpdateUserRequest, type UserResponse } from "../model/user-model";
import { validateWithSchema } from "../validation/validate-with-schema.ts";
import { UserValidation } from "../validation/user-validation.ts";

export const userController = new Hono<{ Variables: ApplicationVariables }>();

userController.delete(
	"/logout",

	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const user: User = c.get("user");
		log.debug(JSON.stringify(user));
		await AuthService.logout(user);
		return c.body(null, 204);
	}
);

userController.get("/current", async (c) => {
	const user: User = c.get("user");

	const response: UserResponse = await UserService.get(user);

	return c.json({
		data: response,
	});
});

userController.patch(
	"/current",
	validateWithSchema("json", UserValidation.UPDATE),
	async (c) => {
		const user: User = c.get("user");
		const request: UpdateUserRequest = await c.req.json();
		const response: UserResponse = await UserService.update(user, request);
		return c.json({
			data: response,
		});
	}
);

import { Hono } from "hono";
import { validator } from "hono/validator";
import type { User } from "@prisma/client";
import { UserService } from "../service/user-service";
import { AuthService } from "../service/auth-service.ts";
import type { ApplicationVariables } from "../model/app-model";
import { type UpdateUserRequest, type UserResponse } from "../model/user-model";

export const userController = new Hono<{ Variables: ApplicationVariables }>();

userController.delete(
	"/logout",

	validator("json", (value) => {
		return value;
	}),

	async (c) => {
		const user: User = c.get("user");
		console.warn(user);
		await AuthService.logout(user);
		return c.body(null, 204);
	}
);

userController.get("/users/current", async (c) => {
	const user: User = c.get("user");

	const response: UserResponse = await UserService.get(user);

	return c.json({
		data: response,
	});
});

userController.patch("/users/current", async (c) => {
	const user: User = c.get("user");

	const request: UpdateUserRequest = await c.req.json();

	const response = await UserService.update(user, request);

	return c.json({
		data: response,
	});
});

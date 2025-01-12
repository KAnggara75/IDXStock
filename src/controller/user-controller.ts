import { Hono } from "hono";
import { log } from "../config/logger";
import type { User } from "@prisma/client";
import { validator } from "hono/validator";
import { UserService } from "../service/user-service";
import { type UpdateUserRequest } from "../model/user-model";
import { authMiddleware } from "../middleware/auth-middleware";
import type { ApplicationVariables } from "../model/app-model";

export const userController = new Hono<{ Variables: ApplicationVariables }>();

userController.post(
	"/users",
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

userController.post(
	"/users/login",

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

userController.use(authMiddleware);

userController.get("/users/current", async (c) => {
	const user: User = c.get("user");

	return c.json({
		data: user,
	});
});

userController.patch("/users/current", async (c) => {
	const user: User = c.get("user") as User;
	const request = (await c.req.json()) as UpdateUserRequest;

	const response = await UserService.update(user, request);

	return c.json({
		data: response,
	});
});

userController.delete("/users/current", async (c) => {
	const user = c.get("user") as User;

	const response = await UserService.logout(user);

	return c.json({
		data: response,
	});
});

import { Hono } from "hono";
import type { User } from "@prisma/client";
import { UserService } from "../service/user-service";
import { type UpdateUserRequest } from "../model/user-model";
import type { ApplicationVariables } from "../model/app-model";

export const userController = new Hono<{ Variables: ApplicationVariables }>();

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

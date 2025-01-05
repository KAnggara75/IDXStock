import { Hono } from "hono";
import {
	type LoginUserRequest,
	type RegisterUserRequest,
	toUserResponse,
	type UpdateUserRequest,
} from "../model/user-model";
import { UserService } from "../service/user-service";
import type { ApplicationVariables } from "../model/app-model";
import type { User } from "@prisma/client";
import { authMiddleware } from "../middleware/auth-middleware";
import { log } from "../config/logger";
import { HTTPException } from "hono/http-exception";

export const userController = new Hono<{ Variables: ApplicationVariables }>();

userController.post("/users", async (c) => {
	const text = await c.req.text();
	let request: RegisterUserRequest;

	try {
		request = JSON.parse(text) as RegisterUserRequest;
	} catch {
		log.error("Invalid json request " + text);
		throw new HTTPException(400, {
			message: "invalid json",
		});
	}

	const response = await UserService.register(request);
	log.info("Registering user", response);

	return c.json({
		data: response,
	});
});

userController.post("/users/login", async (c) => {
	const request = (await c.req.json()) as LoginUserRequest;

	const response = await UserService.login(request);

	return c.json({
		data: response,
	});
});

userController.use(authMiddleware);

userController.get("/users/current", async (c) => {
	const user = c.get("user") as User;

	return c.json({
		data: toUserResponse(user),
	});
});

userController.patch("/users/current", async (c) => {
	const user = c.get("user") as User;
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

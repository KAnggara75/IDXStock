import type { MiddlewareHandler } from "hono";
import { JwtHelper } from "../helpers/jwt-helper";
import { prismaClient } from "../config/database";
import type { UserJwt } from "../model/user-model";

export const authMiddleware: MiddlewareHandler = async (c, next) => {
	const authHeader = c.req.header("Authorization");

	if (!authHeader) {
		return c.json({ errors: "Authorization header is missing" }, 401);
	}

	if (!authHeader.startsWith("Bearer ")) {
		return c.json({ errors: "Authorization header is malformed" }, 401);
	}

	const token = authHeader.split(" ")[1];

	const jwtPayload: UserJwt = await JwtHelper.jwtVerivy(token);

	const user: number = await prismaClient.user.count({
		where: {
			username: jwtPayload.username,
		},
	});

	if (user != 1) {
		return c.json({ errors: "Token has been revoked" }, 401);
	}

	c.set("user", jwtPayload);

	await next();
};

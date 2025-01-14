import type { MiddlewareHandler } from "hono";
import { JwtHelper } from "../helpers/jwt-helper";
import { prismaClient } from "../config/database";
import type { UserJwt } from "../model/user-model";
import { HTTPException } from "hono/http-exception";
import { AuthValidation } from "../validation/auth-validation";

export const authMiddleware: MiddlewareHandler = async (c, next) => {
	const authHeader = c.req.header("Authorization");

	if (!authHeader) {
		return c.json({ errors: "Authorization header is missing" }, 401);
	}

	if (!authHeader.startsWith("Bearer ")) {
		return c.json({ errors: "Authorization header is malformed" }, 401);
	}

	const token = authHeader.split(" ")[1];

	const result = AuthValidation.TOKEN.safeParse(token);

	if (result.error) {
		throw new HTTPException(401, {
			message: "Unauthorized",
		});
	}

	const jwtPayload: UserJwt = await JwtHelper.jwtVerivy(token);

	const user: number = await prismaClient.user.count({
		where: {
			username: jwtPayload.username,
			logoutAt: {
				lte: jwtPayload.iat,
			},
		},
	});

	if (user != 1) {
		return c.json({ errors: "Token has been revoked" }, 401);
	}

	c.set("user", jwtPayload);

	await next();
};

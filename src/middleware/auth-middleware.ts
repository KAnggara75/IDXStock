import { log } from "../config/logger";
import type { MiddlewareHandler } from "hono";
import { JwtHelper } from "../helpers/jwt-helper";
import type { UserJwt } from "../model/user-model";
import { AuthValidation } from "../validation/auth-validation";
import { AuthRepository } from "../repository/auth-repo.ts";

export const authMiddleware: MiddlewareHandler = async (c, next) => {
	console.warn("authMiddleware");

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
		return c.json({ errors: JSON.parse(result.error.message) }, 400);
	}

	const jwtPayload: UserJwt = await JwtHelper.jwtVerify(token);

	log.debug(`jwtPayload ${JSON.stringify(jwtPayload)}`);

	const user: number = await AuthRepository.checkUser(jwtPayload);

	if (user != 1) {
		return c.json({ errors: "Invalid Token" }, 401);
	}

	c.set("user", jwtPayload);

	await next();
};

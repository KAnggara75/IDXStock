import { log } from "../config/logger";
import type { MiddlewareHandler } from "hono";
import { JwtHelper } from "../helpers/jwt-helper";
import type { UserJwt } from "../model/user-model";
import { HTTPException } from "hono/http-exception";
import { AuthRepository } from "../repository/auth-repo.ts";
import { AuthValidation } from "../validation/auth-validation";
import { type CustomError, toErrorDetail } from "../model/errors-model.ts";

export const authMiddleware: MiddlewareHandler = async (c, next) => {
	log.warn("authMiddleware");
	let errorPayload: CustomError[];

	const authHeader = c.req.header("Authorization");

	if (!authHeader) {
		errorPayload = await toErrorDetail(
			"unauthorized",
			"Authorization header is missing",
			[],
			"jwt"
		);
		throw new HTTPException(401, {
			message: JSON.stringify(errorPayload),
		});
	}
	log.warn("authHeader exists");

	if (!authHeader.startsWith("Bearer ")) {
		errorPayload = await toErrorDetail(
			"unauthorized",
			"Authorization header is malformed",
			[],
			"jwt"
		);
		throw new HTTPException(400, {
			message: JSON.stringify(errorPayload),
		});
	}
	log.warn("authHeader startsWith Bearer");

	const token = authHeader.split(" ")[1];
	const result = AuthValidation.TOKEN.safeParse(token);

	if (result.error) {
		return c.json({ errors: JSON.parse(result.error.message) }, 400);
	}

	let jwtPayload: UserJwt;
	try {
		jwtPayload = await JwtHelper.jwtVerify(token);
		log.debug(`jwtPayload ${JSON.stringify(jwtPayload)}`);
	} catch {
		errorPayload = await toErrorDetail(
			"invalid_or_expired_token",
			"Token is invalid or has expired",
			[],
			"jwt"
		);
		throw new HTTPException(401, {
			message: JSON.stringify(errorPayload),
		});
	}

	const user: number = await AuthRepository.checkUser(jwtPayload);

	if (user !== 1) {
		errorPayload = await toErrorDetail(
			"unauthorized_user",
			"Invalid Token: User not authorized",
			[],
			"jwt"
		);
		throw new HTTPException(403, {
			message: JSON.stringify(errorPayload),
		});
	}

	c.set("user", jwtPayload);
	await next();
};

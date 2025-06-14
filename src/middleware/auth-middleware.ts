import { log } from "../config/logger";
import type { MiddlewareHandler } from "hono";
import { JwtHelper } from "../helpers/jwt-helper";
import type { UserJwt } from "../model/user-model";
import { HTTPException } from "hono/http-exception";
import { AuthValidation } from "../validation/auth-validation";
import { type CustomError, toErrorDetail } from "../model/errors-model";
import { RedisService } from "../config/redis";
import { AuthRepository } from "../repository/auth-repository.ts";

export const authMiddleware: MiddlewareHandler = async (c, next) => {
	log.debug("authMiddleware");
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
	log.debug("authHeader exists");

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
	log.debug("authHeader startsWith Bearer");

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
	const redis = new RedisService();
	let logoutAt = await redis.getLogoutAt(jwtPayload.username);
	if (logoutAt === 0) {
		log.debug(`logoutAt is 0, checking from db ${jwtPayload.username}`);
		logoutAt = await AuthRepository.getLogoutAt(jwtPayload);

		if (logoutAt !== 0) {
			await redis.setLogoutAt(jwtPayload.username, logoutAt);
		}
	}
	log.debug(`logoutAt ${logoutAt}`);

	const jwtIat: number = Number(jwtPayload.iat ?? 0);

	if (logoutAt >= jwtIat) {
		log.info(`Token for ${jwtPayload.username} is unauthorized`);
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

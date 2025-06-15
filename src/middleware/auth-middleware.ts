import { log } from "../config/logger";
import type { MiddlewareHandler, Next } from "hono";
import type { SafeParseReturnType } from "zod";
import { JwtHelper } from "../helpers/jwt-helper";
import type { UserJwt } from "../model/user-model";
import { HTTPException } from "hono/http-exception";
import { RedisService } from "../service/redis-service.ts";
import { AuthValidation } from "../validation/auth-validation";
import { AuthRepository } from "../repository/auth-repository.ts";
import { type CustomError, toErrorDetail } from "../model/errors-model";

export const authMiddleware: MiddlewareHandler = async (c, next: Next) => {
	log.debug("authMiddleware");
	let errorPayload: CustomError[];

	const authHeader: string | undefined = c.req.header("Authorization");

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

	const token: string = authHeader.split(" ")[1];
	const result: SafeParseReturnType<unknown, unknown> =
		AuthValidation.TOKEN.safeParse(token);

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
	const username: string = jwtPayload.username;
	let logoutAt: number = await RedisService.getLogoutAt(username);
	if (logoutAt === 0) {
		log.debug(`${username} logoutAt is nil, checking from db ${username}`);
		logoutAt = await AuthRepository.getLogoutAt(jwtPayload);
		await RedisService.setLogoutAt(username, logoutAt);
	}

	log.warn(`${username} logoutAt ${logoutAt}`);

	const jwtIat: number = Number(jwtPayload.iat ?? 0);
	log.warn(`${username} issuedAt ${jwtIat}`);

	if (logoutAt >= jwtIat) {
		log.info(`Token for ${username} is unauthorized`);
		errorPayload = await toErrorDetail(
			"unauthorized_user",
			"Invalid Token: User not authorized",
			[],
			"jwt"
		);
		throw new HTTPException(401, {
			message: JSON.stringify(errorPayload),
		});
	}

	c.set("user", jwtPayload);
	await next();
};

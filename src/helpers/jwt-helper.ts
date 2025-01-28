import { log } from "../config/logger";
import { sign, verify } from "hono/jwt";
import type { User } from "@prisma/client";
import { DateUtils } from "../utils/dateUtils";
import { HTTPException } from "hono/http-exception";
import type { UserJwt } from "../model/user-model";
import {
	JwtAlgorithmNotImplemented,
	JwtTokenExpired,
	JwtTokenInvalid,
	JwtTokenIssuedAt,
	JwtTokenNotBefore,
	JwtTokenSignatureMismatched,
} from "hono/utils/jwt/types";

export class JwtHelper {
	static async jwtSign(user: User): Promise<string> {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";
		const exp = Bun.env.JWT_TOKEN_EXP ?? "1";

		const payload: UserJwt = {
			name: user.name,
			email: user.email,
			username: user.username,
			exp: DateUtils.expInDays(exp),
			iat: Math.floor(Date.now() / 1000),
		};
		log.debug(`create token for ${user.username} exp in ${exp} days`);
		return await sign(payload, secret);
	}

	static async jwtVerivy(token: string) {
		try {
			const secretKey = Bun.env.JWT_SECRET_KEY ?? "";

			const decodedPayload: UserJwt = (await verify(
				token,
				secretKey
			)) as UserJwt;

			const payload: UserJwt = {
				name: decodedPayload.name,
				email: decodedPayload.email,
				username: decodedPayload.username,
				iat: decodedPayload.iat,
				exp: decodedPayload.exp,
			};

			return payload;
			// eslint-disable-next-line @typescript-eslint/no-explicit-any
		} catch (err: any) {
			log.error(`Authorization error ${err.name}`);
			if (err instanceof JwtTokenExpired) {
				throw new HTTPException(401, {
					message: "Session expired, please log in again",
				});
			} else if (err instanceof JwtAlgorithmNotImplemented) {
				throw new HTTPException(401, {
					message: "Jwt Algorithm Not Implemented",
				});
			} else if (err instanceof JwtTokenInvalid) {
				throw new HTTPException(401, {
					message: "Jwt Token Invalid",
				});
			} else if (err instanceof JwtTokenNotBefore) {
				throw new HTTPException(401, {
					message: "Jwt Token Not Before",
				});
			} else if (err instanceof JwtTokenIssuedAt) {
				throw new HTTPException(401, {
					message: "Jwt Token Issued At",
				});
			} else if (err instanceof JwtTokenSignatureMismatched) {
				throw new HTTPException(401, {
					message: "Jwt Token Signature Mismatched",
				});
			} else {
				throw new HTTPException(401, {
					message: "Invalid Authorization",
				});
			}
		}
	}
}

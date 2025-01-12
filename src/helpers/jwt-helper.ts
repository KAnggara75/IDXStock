import { sign, verify } from "hono/jwt";
import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import type { UserJwt, UserResponse } from "../model/user-model";
import { log } from "../config/logger";
import { DateUtils } from "../utils/dateUtils";
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

		log.info(exp);

		const payload = {
			email: user.email,
			name: user.name,
			username: user.username,
			exp: DateUtils.expInDays(exp),
		};
		return await sign(payload, secret);
	}

	static async jwtVerivy(token: string) {
		try {
			const secretKey = Bun.env.JWT_SECRET_KEY ?? "";

			const decodedPayload: UserJwt = (await verify(
				token,
				secretKey
			)) as UserJwt;

			const payload: UserResponse = {
				name: decodedPayload.name,
				email: decodedPayload.email,
				username: decodedPayload.username,
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

import { describe, expect, it } from "bun:test";
import { JwtHelper } from "../src/helpers/jwt-helper";
import type { UserJwt } from "../src/model/user-model";
import { sign } from "hono/jwt";
import { HTTPException } from "hono/http-exception";
import type { User } from "@prisma/client";

describe("JWT Helper", () => {
	it("It success create Token", async () => {
		const payload = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
		} as User;

		try {
			const token: string = await JwtHelper.jwtSign(payload);
			await JwtHelper.jwtVerify(token);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("expired");
			}
		}
	});

	it("It expired", async () => {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";

		const payload: UserJwt = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
			exp: Math.floor(Date.now() / 1000) - 3600,
			iat: Math.floor(Date.now() / 1000),
		};

		try {
			const token: string = await sign(payload, secret);
			await JwtHelper.jwtVerify(token);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("expired");
			}
		}
	});

	it("It Issued At", async () => {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";

		const payload: UserJwt = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
			exp: Math.floor(Date.now() / 1000) + 3600,
			iat: Math.floor(Date.now() / 1000) + 3600,
		};

		try {
			const token: string = await sign(payload, secret);
			await JwtHelper.jwtVerify(token);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("Issued At");
			}
		}
	});

	it("It Jwt Token Not Before", async () => {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";

		const payload: UserJwt = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
			exp: Math.floor(Date.now() / 1000) + 3600,
			nbf: Math.floor(Date.now() / 1000) + 3600,
		};

		try {
			const token: string = await sign(payload, secret);
			await JwtHelper.jwtVerify(token);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("Jwt Token Not Before");
			}
		}
	});

	it("It Jwt Token Signature Mismatched", async () => {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";

		const payload: UserJwt = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
			exp: Math.floor(Date.now() / 1000) + 3600,
		};

		try {
			const token: string = await sign(payload, secret);
			await JwtHelper.jwtVerify(token + "kjnk");
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("Jwt Token Signature Mismatched");
			}
		}
	});

	it("It Invalid Authorization", async () => {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";

		const payload: UserJwt = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
			exp: Math.floor(Date.now() / 1000) + 3600,
		};

		try {
			const token: string = await sign(payload, secret);

			const new_token = token.split(".");

			await JwtHelper.jwtVerify(
				"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzVTMiJ9." +
					new_token[1] +
					"." +
					new_token[2]
			);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("Invalid Authorization");
			}
		}
	});

	it("It Jwt Token Invalid", async () => {
		const secret = Bun.env.JWT_SECRET_KEY ?? "";

		const payload: UserJwt = {
			name: "unitTest",
			username: "username",
			email: "email@gmail.com",
			exp: Math.floor(Date.now() / 1000) + 3600,
		};

		try {
			const token: string = await sign(payload, secret);

			const new_token = token.split(".");

			await JwtHelper.jwtVerify(
				new_token[0] + "." + new_token[1] + "jojo." + new_token[2]
			);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(401);
				expect(err.message).toInclude("Jwt Token Invalid");
			}
		}
	});
});

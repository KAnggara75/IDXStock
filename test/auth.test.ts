import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { log as logger } from "../src/config/logger";
import { RedisTest, UserTest } from "./test-util";
import { app } from "../src/app";
import type { User } from "@prisma/client";
import { JwtHelper } from "../src/helpers/jwt-helper";

describe("POST /api/auth/register", (): void => {
	afterEach(async (): Promise<void> => {
		await UserTest.delete();
		await RedisTest.delete();
	});

	it("should reject register new user if request is invalid", async (): Promise<void> => {
		const response: Response = await app.request("/api/auth/register", {
			method: "post",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: "",
				password: "",
				name: "",
				email: "",
			}),
		});

		const body = await response.json();
		logger.debug(JSON.stringify(body));

		expect(response.status).toBe(400);
		expect(body.errors).toBeDefined();
	});

	it("should reject register new user if request is invalid 2", async (): Promise<void> => {
		const response: Response = await app.request("/api/auth/register", {
			method: "post",
			headers: { "Content-Type": "application/json" },
		});

		const body = await response.json();
		logger.debug(JSON.stringify(body));

		expect(response.status).toBe(400);
		expect(body.errors).toBeDefined();
	});

	it("should reject register new user if username already exists", async (): Promise<void> => {
		await UserTest.create();

		const response: Response = await app.request("/api/auth/register", {
			method: "post",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: "test",
				password: "test",
				name: "test",
				email: "test@gmail.com",
			}),
		});

		const body = await response.json();
		logger.debug(JSON.stringify(body));

		expect(response.status).toBe(400);
		expect(body.errors).toBeDefined();
	});

	it("should register new user success", async (): Promise<void> => {
		const response: Response = await app.request("/api/auth/register", {
			method: "post",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: "test",
				password: "test",
				name: "test",
				email: "test@gmail.com",
			}),
		});

		const body = await response.json();
		logger.debug(JSON.stringify(body));

		expect(response.status).toBe(200);
		expect(body.data).toBeDefined();
		expect(body.data.username).toBe("test");
		expect(body.data.name).toBe("test");
	});
});

describe("POST /api/login", (): void => {
	beforeEach(async (): Promise<void> => {
		await UserTest.create();
	});

	afterEach(async (): Promise<void> => {
		await UserTest.delete();
		await RedisTest.delete();
	});

	it("should be able to login", async (): Promise<void> => {
		const response: Response = await app.request("/api/auth/login", {
			method: "post",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: "test",
				password: "test",
			}),
		});

		expect(response.status).toBe(200);

		const body = await response.json();
		expect(body.data.token).toBeDefined();
	});

	it("should be rejected if username is wrong", async (): Promise<void> => {
		const response: Response = await app.request("/api/auth/login", {
			method: "post",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: "salah",
				password: "test",
			}),
		});

		expect(response.status).toBe(401);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});

	it("should be rejected if password is wrong", async (): Promise<void> => {
		const response: Response = await app.request("/api/auth/login", {
			method: "post",
			headers: { "Content-Type": "application/json" },
			body: JSON.stringify({
				username: "test",
				password: "salah",
			}),
		});

		expect(response.status).toBe(401);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});
});

describe("DELETE /api/users/logout", (): void => {
	let token: string = "";

	beforeEach(async (): Promise<void> => {
		const user: User = await UserTest.create();
		const now: number = Math.floor(Date.now() / 1000);
		token = await JwtHelper.jwtSign(user, now - 10);
	});

	afterEach(async (): Promise<void> => {
		await UserTest.delete();
		await RedisTest.delete();
	});

	it("should failed logout cause user not registered", async (): Promise<void> => {
		await UserTest.delete();
		const logout: Response = await app.request("/api/users/logout", {
			method: "DELETE",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});
		expect(logout.status).toBe(403);
	});

	it("should failed logout cause invalid request", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/logout", {
			method: "delete",
		});

		expect(response.status).toBe(401);
		const body = await response.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized");
		expect(body.errors[0].message).toBe("Authorization header is missing");
	});

	it("should failed logout cause invalid token", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/logout", {
			method: "delete",
			headers: {
				Authorization: `Bearer ${token}a`,
			},
		});

		expect(response.status).toBe(401);
		const body = await response.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("invalid_or_expired_token");
		expect(body.errors[0].message).toBe("Token is invalid or has expired");
	});

	it("should failed logout cause invalid header auth format", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/logout", {
			method: "delete",
			headers: {
				Authorization: `${token}`,
			},
		});

		expect(response.status).toBe(400);
		const body = await response.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized");
		expect(body.errors[0].message).toBe("Authorization header is malformed");
	});

	it("should be success logout", async (): Promise<void> => {
		const logout: Response = await app.request("/api/users/logout", {
			method: "DELETE",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});
		expect(logout.status).toBe(204);

		const getUser: Response = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});
		expect(getUser.status).toBe(401);
		const body = await getUser.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized_user");
		expect(body.errors[0].message).toBe("Invalid Token: User not authorized");
	});

	it("test with invalid jwt zod validation", async (): Promise<void> => {
		const logout: Response = await app.request("/api/users/logout", {
			method: "DELETE",
			headers: {
				Authorization: "Bearer test",
			},
		});

		const body = await logout.json();
		expect(logout.status).toBe(400);
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("invalid_string");
		expect(body.errors[0].message).toBe("Invalid jwt");
	});
});

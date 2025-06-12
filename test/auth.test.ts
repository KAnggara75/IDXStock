import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { log as logger } from "../src/config/logger";
import { UserTest } from "./test-util";
import { app } from "../src";
import type { User } from "@prisma/client";
import { JwtHelper } from "../src/helpers/jwt-helper";

describe("POST /api/auth/register", () => {
	afterEach(async () => {
		await UserTest.delete();
	});

	it("should reject register new user if request is invalid", async () => {
		const response = await app.request("/api/auth/register", {
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

	it("should reject register new user if request is invalid 2", async () => {
		const response = await app.request("/api/auth/register", {
			method: "post",
			headers: { "Content-Type": "application/json" },
		});

		const body = await response.json();
		logger.debug(JSON.stringify(body));

		expect(response.status).toBe(400);
		expect(body.errors).toBeDefined();
	});

	it("should reject register new user if username already exists", async () => {
		await UserTest.create();

		const response = await app.request("/api/auth/register", {
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

	it("should register new user success", async () => {
		const response = await app.request("/api/auth/register", {
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

describe("POST /api/login", () => {
	beforeEach(async () => {
		await UserTest.create();
	});

	afterEach(async () => {
		await UserTest.delete();
	});

	it("should be able to login", async () => {
		const response = await app.request("/api/auth/login", {
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

	it("should be rejected if username is wrong", async () => {
		const response = await app.request("/api/auth/login", {
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

	it("should be rejected if password is wrong", async () => {
		const response = await app.request("/api/auth/login", {
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

describe("DELETE /api/users/logout", () => {
	let token: string = "";

	beforeEach(async () => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async () => {
		await UserTest.delete();
	});

	it("should failed logout cause user not registered", async () => {
		const logout = await app.request("/api/users/logout", {
			method: "DELETE",
			headers: {
				Authorization:
					"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoibm90Zm91bmQiLCJlbWFpbCI6Im5vdGZvdW5kQGdtYWlsLmNvbSIsInVzZXJuYW1lIjoibm90Zm91bmQiLCJleHAiOjE3NTIzMzQ2MDcsImlhdCI6MTc0OTc0MjYwN30.Uif93W-NnzmprpnFKXUF4_HwalWNySRE9R-yh6I_A_0",
			},
		});
		expect(logout.status).toBe(403);
	});

	it("should failed logout cause invalid request", async () => {
		const response = await app.request("/api/users/logout", {
			method: "delete",
		});

		expect(response.status).toBe(401);
		const body = await response.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized");
		expect(body.errors[0].message).toBe("Authorization header is missing");
	});

	it("should failed logout cause invalid token", async () => {
		const response = await app.request("/api/users/logout", {
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

	it("should failed logout cause invalid header auth format", async () => {
		const response = await app.request("/api/users/logout", {
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

	it("should be success logout", async () => {
		const logout = await app.request("/api/users/logout", {
			method: "DELETE",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});
		expect(logout.status).toBe(204);

		const getUser = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: `Bearer ${token}`,
			},
		});
		expect(getUser.status).toBe(403);
		const body = await getUser.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized_user");
		expect(body.errors[0].message).toBe("Invalid Token: User not authorized");
	});

	it("test with invalid jwt zod validation", async () => {
		const logout = await app.request("/api/users/logout", {
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

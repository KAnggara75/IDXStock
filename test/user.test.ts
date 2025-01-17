import { describe, it, expect, afterEach, beforeEach } from "bun:test";
import { log as logger } from "../src/config/logger";
import { app } from "../src";
import { UserTest } from "./test-util";
import type { User } from "@prisma/client";
import { JwtHelper } from "../src/helpers/jwt-helper";

describe("GET /api/users/current", () => {
	let token: string = "";

	beforeEach(async () => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async () => {
		await UserTest.delete();
	});

	it("should be able to get user", async () => {
		const response = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: `Bearer ${token}`,
				"Content-Type": "application/json",
			},
		});

		expect(response.status).toBe(200);

		const body = await response.json();
		expect(body.data).toBeDefined();
		expect(body.data.username).toBe("test");
		expect(body.data.name).toBe("test");
	});

	it("should not be able to get user if token is invalid", async () => {
		const response = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: "salah",
				"Content-Type": "application/json",
			},
		});

		expect(response.status).toBe(401);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});

	it("should not be able to get user if there is no Authorization header", async () => {
		const response = await app.request("/api/users/current", {
			headers: { "Content-Type": "application/json" },
			method: "get",
		});

		expect(response.status).toBe(401);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});
});

describe("PATCH /api/users/current", () => {
	let token: string = "";

	beforeEach(async () => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async () => {
		await UserTest.delete();
	});

	it("should be rejected if request is invalid", async () => {
		const response = await app.request("/api/users/current", {
			method: "patch",
			headers: {
				Authorization: `Bearer ${token}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				name: "",
				password: "",
			}),
		});

		expect(response.status).toBe(400);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});

	it("should be able to update name", async () => {
		const response = await app.request("/api/users/current", {
			method: "patch",
			headers: {
				Authorization: `Bearer ${token}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				name: "IDScript",
			}),
		});

		expect(response.status).toBe(200);

		const body = await response.json();
		logger.debug(JSON.stringify(body));
		expect(body.data).toBeDefined();
		expect(body.data.name).toBe("IDScript");
	});

	it("should be able to update password", async () => {
		let response = await app.request("/api/users/current", {
			method: "patch",
			headers: {
				Authorization: `Bearer ${token}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				password: "passwordBaru",
			}),
		});

		expect(response.status).toBe(200);

		const body = await response.json();
		logger.debug(JSON.stringify(body));
		expect(body.data).toBeDefined();
		expect(body.data.name).toBe("test");

		response = await app.request("/api/login", {
			headers: { "Content-Type": "application/json" },
			method: "post",
			body: JSON.stringify({
				username: "test",
				password: "passwordBaru",
			}),
		});

		expect(response.status).toBe(200);
	});
});

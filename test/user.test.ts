import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { log as logger } from "../src/config/logger";
import { app } from "../src";
import { RedisTest, UserTest } from "./test-util";
import type { User } from "@prisma/client";
import { JwtHelper } from "../src/helpers/jwt-helper";

describe("GET /api/users/current", () => {
	let token: string = "";

	beforeEach(async (): Promise<void> => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async (): Promise<void> => {
		await UserTest.delete();
		await RedisTest.delete();
	});

	it("should be able to get user", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
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

	it("should not be able to get user if token is Invalid jwt", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: "Bearer salah",
				"Content-Type": "application/json",
			},
		});

		expect(response.status).toBe(400);

		const body = await response.json();
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("invalid_string");
		expect(body.errors[0].message).toBe("Invalid jwt");
	});

	it("should not be able to get user if token is Invalid", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: "salah",
				"Content-Type": "application/json",
			},
		});

		const body = await response.json();
		expect(response.status).toBe(400);
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized");
		expect(body.errors[0].message).toBe("Authorization header is malformed");
	});

	it("should not be able to get user if there is no Authorization header", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
			headers: { "Content-Type": "application/json" },
			method: "get",
		});

		expect(response.status).toBe(401);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});
});

describe("PATCH /api/users/current", () => {
	let userToken: string = "";

	beforeEach(async (): Promise<void> => {
		const user: User = await UserTest.create();
		userToken = await JwtHelper.jwtSign(user);
	});

	afterEach(async (): Promise<void> => {
		await UserTest.delete();
		await RedisTest.delete();
	});

	it("should be rejected if request is invalid", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
			method: "patch",
			headers: {
				Authorization: `Bearer ${userToken}`,
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

	it("should be able to update name", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
			method: "patch",
			headers: {
				Authorization: `Bearer ${userToken}`,
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

	it("should be able to update password", async (): Promise<void> => {
		const response: Response = await app.request("/api/users/current", {
			method: "patch",
			headers: {
				Authorization: `Bearer ${userToken}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				password: "passwordBaru",
			}),
		});

		expect(response.status).toBe(200);

		let body = await response.json();
		logger.debug(JSON.stringify(body));
		expect(body.data).toBeDefined();
		expect(body.data.name).toBe("test");

		const userRes: Response = await app.request("/api/users/current", {
			method: "get",
			headers: {
				Authorization: `Bearer ${userToken}`,
				"Content-Type": "application/json",
			},
		});

		body = await userRes.json();
		expect(userRes.status).toBe(401);
		expect(body.errors).toBeDefined();
		expect(body.errors[0].validation).toBe("jwt");
		expect(body.errors[0].code).toBe("unauthorized_user");
		expect(body.errors[0].message).toBe("Invalid Token: User not authorized");

		// response = await app.request("/api/users/login", {
		// 	headers: { "Content-Type": "application/json" },
		// 	method: "post",
		// 	body: JSON.stringify({
		// 		username: "test",
		// 		password: "passwordBaru",
		// 	}),
		// });
		//
		// expect(response.status).toBe(200);
	});
});

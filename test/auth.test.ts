import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { log as logger } from "../src/config/logger";
import { UserTest } from "./test-util";
import { app } from "../src";

describe("POST /api/register", () => {
	afterEach(async () => {
		await UserTest.delete();
	});

	it("should reject register new user if request is invalid", async () => {
		const response = await app.request("/api/register", {
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
		const response = await app.request("/api/register", {
			method: "post",
		});

		const body = await response.json();
		logger.debug(JSON.stringify(body));

		expect(response.status).toBe(400);
		expect(body.errors).toBeDefined();
	});

	it("should reject register new user if username already exists", async () => {
		await UserTest.create();

		const response = await app.request("/api/register", {
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
		const response = await app.request("/api/register", {
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
		logger.info(JSON.stringify(body));

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
		const response = await app.request("/api/login", {
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
		const response = await app.request("/api/login", {
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
		const response = await app.request("/api/login", {
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

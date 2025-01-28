import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import type { User } from "@prisma/client";
import { DailyTest, UserTest } from "./test-util.ts";
import { JwtHelper } from "../src/helpers/jwt-helper.ts";
import { app } from "../src";

describe("POST /api/stocks/idx", () => {
	let token: string = "";

	beforeEach(async () => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async () => {
		await UserTest.delete();
	});

	it("should failed cause invalid token", async () => {
		const response = await app.request("api/stocks/idx", {
			method: "post",
		});

		expect(response.status).toBe(401);
	});

	it("should failed cause invalid request", async () => {
		const response = await app.request("api/stocks/idx", {
			method: "post",
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
	});
});

describe("PATCH /api/stocks/google", () => {
	let token: string = "";

	beforeEach(async () => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async () => {
		await DailyTest.delete();
		await UserTest.delete();
	});

	it("should failed cause invalid token", async () => {
		const response = await app.request("api/stocks/google", {
			method: "patch",
		});

		expect(response.status).toBe(401);
	});

	it("should be success update Stock daily", async () => {
		const response = await app.request("api/stocks/google", {
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

		expect(response.status).toBe(200);
	});
});

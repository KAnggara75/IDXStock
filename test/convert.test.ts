import { afterEach, beforeEach, describe, expect, it } from "bun:test";
import { app } from "../src";
import { UserTest } from "./test-util";
import type { User } from "@prisma/client";
import { JwtHelper } from "../src/helpers/jwt-helper";
import * as fs from "node:fs";

describe("POST /api/convert", () => {
	let token: string = "";

	beforeEach(async () => {
		const user: User = await UserTest.create();
		token = await JwtHelper.jwtSign(user);
	});

	afterEach(async () => {
		await UserTest.delete();
	});

	it("should be able convert xlsx to json", async () => {
		const formData = new FormData();
		formData.append(
			"file",
			fs.createReadStream("./test/Stock Summary-12345678.xlsx")
		);

		const response = await app.request("/api/convert", {
			method: "post",
			headers: {
				Authorization: `Bearer ${token}`,
			},
			body: formData,
		});

		const body = await response.json();

		expect(body.errors).toBeDefined();
		expect(body.errors[0].code).toBeDefined();
		expect(body.errors[0].message).toBeDefined();
		expect(body.errors[0].fatal).toBe(true);
		expect(response.status).toBe(400);
	});
});

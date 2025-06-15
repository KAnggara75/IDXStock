import { app } from "../src/app.ts";
import { describe, expect, it } from "bun:test";

describe("POST", (): void => {
	it("URL Not Found", async (): Promise<void> => {
		const url: `${string}-${string}-${string}-${string}-${string}` =
			crypto.randomUUID();
		const response: Response = await app.request("/" + url, {
			method: "get",
		});

		expect(response.status).toBe(404);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});

	it("App Version", async (): Promise<void> => {
		const response: Response = await app.request("/api/version", {
			method: "get",
		});

		expect(response.status).toBe(200);

		const body = await response.json();
		expect(body.message).toBeDefined();
		expect(body.version).toBeDefined();
		expect(body.errors).toBeUndefined();
		expect(body.stability).toBeDefined();
	});
});

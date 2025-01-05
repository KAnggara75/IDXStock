import { app } from "../src";
import { expect, describe, it } from "bun:test";

describe("POST", () => {
	it("URL Not Found", async () => {
		const url = crypto.randomUUID();
		const response = await app.request("/" + url, {
			method: "get",
			body: JSON.stringify({
				first_name: "OK",
			}),
		});

		expect(response.status).toBe(404);

		const body = await response.json();
		expect(body.errors).toBeDefined();
	});

	it("App Version", async () => {
		const response = await app.request("/api", {
			method: "get",
			body: JSON.stringify({
				first_name: "OK",
			}),
		});

		expect(response.status).toBe(200);

		const body = await response.json();
		expect(body.message).toBeDefined();
		expect(body.version).toBeDefined();
		expect(body.stability).toBeDefined();
		expect(body.errors).toBeUndefined();
	});
});

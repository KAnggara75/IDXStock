import { describe, expect, it } from "bun:test";
import {
	type Error,
	type GoogleFinance,
	type StockModel,
	toGoogleFinance,
} from "../src/model/googleFinance-model";
import { HTTPException } from "hono/http-exception";

describe("Google Finance", () => {
	const value1: string[] = [
		"code",
		"name",
		"price",
		"priceopen",
		"high",
		"low",
		"volume",
		"marketcap",
		"tradetime",
		"volumeavg",
		"pe",
		"eps",
		"high52",
		"low52",
		"change",
		"changepct",
		"closeyest",
		"shares",
	];

	it("should Bad data GOOGLE_FINANCE_FORMAT length", () => {
		const value1: string[] = ["a", "a", "a", "a", "a", "a", "a", "a"];
		const value: string[][] = [value1, value1];

		const testData: GoogleFinance = {
			values: value,
		};

		try {
			toGoogleFinance(testData);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(422);
				expect(err.message).toBe("Bad data GOOGLE_FINANCE_FORMAT length");
			}
		}
	});

	it("should Bad data GOOGLE_FINANCE_FORMAT header name", () => {
		const value2: string[] = [
			"code",
			"name",
			"price",
			"priceopen",
			"high",
			"low",
			"volume",
			"marketcap",
			"tradetime",
			"volumeavg",
			"pe",
			"eps",
			"high52",
			"low52",
			"change",
			"changepct",
			"closeyest",
			"sharess",
		];
		const value: string[][] = [value2, value1];

		const testData: GoogleFinance = {
			values: value,
		};

		try {
			toGoogleFinance(testData);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(422);
				expect(err.message).toBe("Bad data GOOGLE_FINANCE_FORMAT header name");
			}
		}
	});

	it("should error", () => {
		const error: Error = {
			code: 404,
			message: "Not Found",
			status: "Not Found",
		};
		const testData: GoogleFinance = {
			error: error,
		};

		try {
			toGoogleFinance(testData);
		} catch (err) {
			if (err instanceof HTTPException) {
				expect(err.status).toBe(422);
				expect(err.message).toBe("Bad data GoogleFinance");
			}
		}
	});

	it("should error", () => {
		const value: string[][] = [value1, value1];

		const testData: GoogleFinance = {
			values: value,
		};

		const data: StockModel[] = toGoogleFinance(testData);

		expect(data[0].code).toBe("code");
		expect(data[0].name).toBe("name");
		expect(data[0].price).toBe(0);
		expect(data[0].priceopen).toBe(0);
		expect(data[0].high).toBe(0);
		expect(data[0].low).toBe(0);
		expect(data[0].volume).toBe(0);
		expect(data[0].marketcap).toBe(0);
		expect(data[0].volumeavg).toBe(0);
		expect(data[0].pe).toBe(0);
		expect(data[0].eps).toBe(0);
		expect(data[0].high52).toBe(0);
		expect(data[0].low52).toBe(0);
		expect(data[0].change).toBe(0);
		expect(data[0].changepct).toBe(0);
		expect(data[0].closeyest).toBe(0);
		expect(data[0].shares).toBe(0);
	});
});

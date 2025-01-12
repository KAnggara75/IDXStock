import { describe, expect, it } from "bun:test";
import { DateUtils } from "../src/utils/dateUtils.ts";

describe("expInDays", () => {
	it("correct days", () => {
		const date = DateUtils.expInDays(30);
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24 * 30;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});

	it("correct number in string", () => {
		const date = DateUtils.expInDays("30");
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24 * 30;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});

	it("invalid number 30a", () => {
		const date = DateUtils.expInDays("30a");
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});

	it("invalid number aabb", () => {
		const date = DateUtils.expInDays("aabb");
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});
});

describe("toDate", () => {
	it("correct days", () => {
		const date = DateUtils.toDate("2025-01-10 16:09:51");
		const expectedDate: Date = new Date(
			new Date("2025-01-10T16:09:51.000Z").getTime() + 7 * 60 * 60 * 1000
		);

		expect(date).toBeDefined();
		expect(date).toEqual(expectedDate);
	});

	it("incorrect days", () => {
		const date = DateUtils.toDate("");
		const expectedDate: Date = new Date();
		expect(date).toBeDefined();
		expect(date).toEqual(expectedDate);
	});
});

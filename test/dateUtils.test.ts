import { describe, expect, it } from "bun:test";
import { DateUtils } from "../src/utils/dateUtils";

describe("expInDays", (): void => {
	it("correct days", (): void => {
		const date: number = DateUtils.expInDays(30);
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24 * 30;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});

	it("correct number in string", (): void => {
		const date: number = DateUtils.expInDays("30");
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24 * 30;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});

	it("invalid number 30a", (): void => {
		const date: number = DateUtils.expInDays("30a");
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});

	it("invalid number aabb", (): void => {
		const date: number = DateUtils.expInDays("aabb");
		const expectedDate: number = Math.floor(Date.now() / 1000) + 3600 * 24;

		expect(date).toBeDefined();
		expect(date).toBe(expectedDate);
		expect(date).toEqual(expectedDate);
	});
});

describe("toDate", (): void => {
	it("correct days", (): void => {
		const date: Date = DateUtils.toDate("2025-01-10 16:09:51");
		const expectedDate: Date = new Date(
			new Date("2025-01-10T16:09:51.000Z").getTime() + 7 * 60 * 60 * 1000
		);

		expect(date).toBeDefined();
		expect(date).toEqual(expectedDate);
	});

	it("incorrect days", (): void => {
		const date: Date = DateUtils.toDate("");
		const expectedDate: Date = new Date();
		expect(date).toBeDefined();
		expect(date).toEqual(expectedDate);
	});
});

describe("toDate2", (): void => {
	it("correct days", (): void => {
		const date: Date = DateUtils.toDate2("17 Jan 2025");
		const expectedDate: Date = new Date(
			new Date("2025-01-17").getTime() + 7 * 60 * 60 * 1000
		);

		expect(date).toBeDefined();
		expect(date).toEqual(expectedDate);
	});

	it("incorrect days", (): void => {
		const date: Date = DateUtils.toDate2("");
		const expectedDate: Date = new Date();
		expect(date).toBeDefined();
		expect(date).toEqual(expectedDate);
	});
});

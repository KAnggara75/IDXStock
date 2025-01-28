import { describe, expect, it } from "bun:test";
import { NumberUtil } from "../src/utils/numberUtils.ts";

describe("NumberUtil to Int", () => {
	it.each([
		[2, 2],
		[-1, -1],
		[3.1, 3],
		[4.9, 4],
		[0.9, 0],
	])("%d => %d", (a: string | number, expected: number) => {
		const result: number = NumberUtil.toInt(a);
		expect(result).toBe(expected);
		expect(result).toBeTypeOf("number");
	});

	it.each([
		["", 0],
		["2", 2],
		["a", 0],
		["-3", -3],
		["0,1", 0],
		["1,2", 1],
		["2,1", 2],
		["1m3", 0],
		["1,m3", 0],
		["1.m3", 0],
		["1.15503E+11", 115503000000],
		["1,15503E+11", 115503000000],
	])("%s => %d", (a: string | number, expected: number) => {
		const result: number = NumberUtil.toInt(a);
		expect(result).toBe(expected);
		expect(result).toBeTypeOf("number");
	});
});

describe("NumberUtil to Float", () => {
	it.each([
		[1, 1.0],
		[2, 2.0],
		[-3, -3.0],
		[3.1, 3.1],
		[4.9, 4.9],
		[0.9, 0.9],
	])("%d => %d", (a: string | number, expected: number) => {
		const result: number = NumberUtil.toFloat(a);
		expect(result).toBe(expected);
		expect(result).toBeTypeOf("number");
	});

	it.each([
		["", 0.0],
		["2", 2.0],
		["a", 0.0],
		["-3", -3.0],
		["0,1", 0.1],
		["1,2", 1.2],
		["2,1", 2.1],
		["1m3", 0.0],
		["1,m3", 0.0],
		["1.m3", 0.0],
		["1.15503E+11", 115503000000.0],
		["1,15503E+11", 115503000000.0],
	])("%s => %d", (a: string | number, expected: number) => {
		const result: number = NumberUtil.toFloat(a);
		expect(result).toBe(expected);
		expect(result).toBeTypeOf("number");
	});
});

describe("NumberUtil to BigInt", () => {
	it.each([
		[2, 2],
		[-1, -1],
		[3.1, 3],
		[4.9, 4],
		[0.9, 0],
	])("%d => %d", (a: number, expected: number) => {
		const result: bigint = NumberUtil.toBigInt(a);
		expect(result).toBe(BigInt(expected));
		expect(result).toBeTypeOf("bigint");
	});
	it.each([
		["", 0],
		["2", 2],
		["a", 0],
		["-3", -3],
		["0,1", 0],
		["1,2", 1],
		["2,1", 2],
		["1m3", 0],
		["1,m3", 0],
		["1.m3", 0],
		["1.15503E+11", 115503000000],
		["1,15503E+11", 115503000000],
	])("%s => %d", (a: string, expected: number) => {
		const result: bigint = NumberUtil.toBigInt(a);
		expect(result).toBe(BigInt(expected));
		expect(result).toBeTypeOf("bigint");
	});
});

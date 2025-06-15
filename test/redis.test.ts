import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { RedisService } from "../src/service/redis-service.ts";

const testKey = "test-key";
const testUser = "test-user";

beforeAll(async (): Promise<void> => {
	await RedisService.del(testKey);
	await RedisService.del(`idx-${testUser}`);
});

afterAll(async (): Promise<void> => {
	await RedisService.del(testKey);
	await RedisService.del(`idx-${testUser}`);
});

describe("RedisService", (): void => {
	it("should set and get a key", async (): Promise<void> => {
		await RedisService.set(testKey, "hello", 10);
		const result: string | null = await RedisService.get(testKey);
		expect(result).toBe("hello");
	});

	it("should delete a key", async (): Promise<void> => {
		await RedisService.set(testKey, "to-delete");
		await RedisService.del(testKey);
		const result: string | null = await RedisService.get(testKey);
		expect(result).toBeNull();
	});

	it("should return true if key exists", async (): Promise<void> => {
		await RedisService.set(testKey, "exists-check");
		const exists: boolean = await RedisService.exists(testKey);
		expect(exists).toBe(true);
	});

	it("should return false if key does not exist", async (): Promise<void> => {
		await RedisService.del(testKey);
		const exists: boolean = await RedisService.exists(testKey);
		expect(exists).toBe(false);
	});

	it("should set and get logoutAt", async (): Promise<void> => {
		const timestamp: number = Math.floor(Date.now() / 1000);
		await RedisService.setLogoutAt(testUser, timestamp);
		const result: number = await RedisService.getLogoutAt(testUser);
		expect(result).toBe(timestamp);
	});

	it("should return 0 for logoutAt if key missing", async (): Promise<void> => {
		await RedisService.del(`idx-${testUser}`);
		const result: number = await RedisService.getLogoutAt(testUser);
		expect(result).toBe(0);
	});
});

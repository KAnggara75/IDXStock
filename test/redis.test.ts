import { afterAll, beforeAll, describe, expect, it } from "bun:test";
import { RedisService } from "../src/config/redis.ts";

const redis = new RedisService();
const testKey = "test-key";
const testUser = "test-user";

beforeAll(async () => {
	await redis.del(testKey);
	await redis.del(`idx-${testUser}`);
});

afterAll(async () => {
	await redis.del(testKey);
	await redis.del(`idx-${testUser}`);
});

describe("RedisService", () => {
	it("should set and get a key", async () => {
		await redis.set(testKey, "hello", 10);
		const result = await redis.get(testKey);
		expect(result).toBe("hello");
	});

	it("should delete a key", async () => {
		await redis.set(testKey, "to-delete");
		await redis.del(testKey);
		const result = await redis.get(testKey);
		expect(result).toBeNull();
	});

	it("should return true if key exists", async () => {
		await redis.set(testKey, "exists-check");
		const exists = await redis.exists(testKey);
		expect(exists).toBe(true);
	});

	it("should return false if key does not exist", async () => {
		await redis.del(testKey);
		const exists = await redis.exists(testKey);
		expect(exists).toBe(false);
	});

	it("should set and get logoutAt", async () => {
		const timestamp = Math.floor(Date.now() / 1000);
		await redis.setLogoutAt(testUser, timestamp);
		const result = await redis.getLogoutAt(testUser);
		expect(result).toBe(timestamp);
	});

	it("should return 0 for logoutAt if key missing", async () => {
		await redis.del(`idx-${testUser}`);
		const result = await redis.getLogoutAt(testUser);
		expect(result).toBe(0);
	});
});

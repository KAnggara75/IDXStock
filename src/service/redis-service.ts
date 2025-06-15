import { redisClient } from "../config/redis";

export class RedisService {
	static async set(
		key: string,
		value: string,
		ttlSeconds?: number
	): Promise<void> {
		if (ttlSeconds) {
			await redisClient.set(key, value, "EX", ttlSeconds);
		} else {
			await redisClient.set(key, value);
		}
	}

	static async get(key: string): Promise<string | null> {
		return redisClient.get(key);
	}

	static async del(key: string): Promise<void> {
		await redisClient.del(key);
	}

	static async exists(key: string): Promise<boolean> {
		return (await redisClient.exists(key)) === 1;
	}

	static async getLogoutAt(username: string): Promise<number> {
		const key = `idx-${username}`;
		const value = await this.get(key);
		return value ? Number(value) : 0;
	}

	static async setLogoutAt(
		username: string,
		timestamp: number,
		ttlSeconds = 60 * 60 * 24 * 7
	): Promise<void> {
		const key = `idx-${username}`;
		await this.set(key, timestamp.toString(), ttlSeconds);
	}
}

import { redisClient } from "../config/redis";
import { log } from "../config/logger.ts";

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
		const value: string | null = await this.get(key);
		return value ? Number(value) : 0;
	}

	static async setLogoutAt(
		username: string,
		timestamp: number,
		ttlSeconds: number = 60 * 60 * 24 * 7
	): Promise<void> {
		const key = `idx-${username}`;
		await this.set(key, timestamp.toString(), ttlSeconds);
		log.debug(`${username} logged out at ${timestamp}`);
	}
}

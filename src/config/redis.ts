import Redis from "ioredis";

export class RedisService {
	private client: Redis;

	constructor() {
		const url = process.env.REDIS_URL || "redis://localhost:6379";
		this.client = new Redis(url);
	}

	async set(key: string, value: string, ttlSeconds?: number): Promise<void> {
		if (ttlSeconds) {
			await this.client.set(key, value, "EX", ttlSeconds);
		} else {
			await this.client.set(key, value);
		}
	}

	async get(key: string): Promise<string | null> {
		return this.client.get(key);
	}

	async del(key: string): Promise<void> {
		await this.client.del(key);
	}

	async exists(key: string): Promise<boolean> {
		return (await this.client.exists(key)) === 1;
	}

	async getLogoutAt(username: string): Promise<number> {
		const key = `idx-${username}`;
		const value = await this.get(key);
		return value ? Number(value) : 0;
	}

	async setLogoutAt(
		username: string,
		timestamp: number,
		ttlSeconds = 60 * 60 * 24 * 7
	): Promise<void> {
		const key = `idx-${username}`;
		await this.set(key, timestamp.toString(), ttlSeconds);
	}
}

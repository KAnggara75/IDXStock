import Redis from "ioredis";
import { log } from "./logger";
import type { Logger } from "winston";

const url = process.env.REDIS_URL || "redis://localhost:6379";
export const redisClient = new Redis(url);

redisClient.on("error", (err: Error): void => {
	log.error("Redis error:", err);
});

redisClient.on("connect", (): Logger => log.debug("Redis CONNECT"));
redisClient.on("ready", (): Logger => log.debug("Redis READY"));

import type { User } from "@prisma/client";
import { prismaClient } from "../src/config/database";
import { RedisService } from "../src/service/redis-service.ts";

export class UserTest {
	static async create(): Promise<User> {
		return prismaClient.user.create({
			data: {
				username: "test",
				name: "test",
				email: "test@gmail.com",
				password: await Bun.password.hash("test", {
					algorithm: "bcrypt",
					cost: 10,
				}),
			},
		});
	}

	static async delete() {
		await prismaClient.user.deleteMany({
			where: {
				username: "test",
			},
		});
	}
}

export class DailyTest {
	static async delete() {
		await prismaClient.daily.deleteMany({
			where: {
				insertBy: "test",
			},
		});
	}
}

export class HistoryTest {
	static async delete() {
		await prismaClient.history.deleteMany({
			where: {
				code: "TEST",
			},
		});
	}
}

export class RedisTest {
	static async delete() {
		await RedisService.del("idx-test");
	}
}

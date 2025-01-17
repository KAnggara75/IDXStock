import type { User } from "@prisma/client";
import { prismaClient } from "../src/config/database";

export class UserTest {
	static async create(): Promise<User> {
		return await prismaClient.user.create({
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

import { prismaClient } from "../config/database.ts";
import type { UserJwt } from "../model/user-model.ts";

export class AuthRepository {
	public constructor() {}
	static async checkUser(jwt: UserJwt): Promise<number> {
		return prismaClient.user.count({
			where: {
				username: jwt.username,
				logoutAt: {
					lt: jwt.iat,
				},
			},
		});
	}
}

import { prismaClient } from "../config/database";
import type { UserJwt } from "../model/user-model";

export class AuthRepository {
	public constructor() {}
	static async getLogoutAt(jwt: UserJwt): Promise<number> {
		const user = await prismaClient.user.findUnique({
			where: { username: jwt.username },
			select: { logoutAt: true },
		});

		return Number(user?.logoutAt ?? 0);
	}
}

import {
	toUserResponse,
	type UpdateUserRequest,
	type UserResponse,
} from "../model/user-model";
import type { User } from "@prisma/client";
import { prismaClient } from "../config/database";
import { HTTPException } from "hono/http-exception";
import { UserValidation } from "../validation/user-validation";
import { AuthValidation } from "../validation/auth-validation.ts";

export class UserService {
	static async get(token: string | undefined | null): Promise<User> {
		const result = AuthValidation.TOKEN.safeParse(token);

		if (result.error) {
			throw new HTTPException(401, {
				message: "Unauthorized",
			});
		}

		token = result.data;

		const user = await prismaClient.user.findFirst({
			where: {
				token: token,
			},
		});

		if (!user) {
			throw new HTTPException(401, {
				message: "Unauthorized",
			});
		}

		return user;
	}

	static async update(
		user: User,
		request: UpdateUserRequest
	): Promise<UserResponse> {
		request = UserValidation.UPDATE.parse(request);

		if (request.name) {
			user.name = request.name;
		}

		if (request.password) {
			user.password = await Bun.password.hash(request.password, {
				algorithm: "bcrypt",
				cost: 10,
			});
		}

		user = await prismaClient.user.update({
			where: {
				username: user.username,
			},
			data: user,
		});

		return toUserResponse(user);
	}

	static async logout(user: User): Promise<boolean> {
		await prismaClient.user.update({
			where: {
				username: user.username,
			},
			data: {
				token: null,
			},
		});

		return true;
	}
}

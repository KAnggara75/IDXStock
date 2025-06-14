import {
	toUserResponse,
	type UpdateUserRequest,
	type UserResponse,
} from "../model/user-model";
import type { User } from "@prisma/client";
import { AuthService } from "./auth-service";
import { prismaClient } from "../config/database";

export class UserService {
	static async get(token: User): Promise<UserResponse> {
		const user: User | null = await prismaClient.user.findFirst({
			where: {
				username: token.username,
			},
		});

		return toUserResponse(user!);
	}

	static async update(
		user: User,
		request: UpdateUserRequest
	): Promise<UserResponse> {
		if (request.password) {
			request.password = await Bun.password.hash(request.password, {
				algorithm: "bcrypt",
				cost: 10,
			});
		}

		user = await prismaClient.user.update({
			where: {
				username: user.username,
			},
			data: request,
		});

		await AuthService.logout(user);

		return toUserResponse(user, true);
	}
}

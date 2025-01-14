import {
	toUserResponse,
	type UserResponse,
	type UpdateUserRequest,
} from "../model/user-model";
import type { User } from "@prisma/client";
import { prismaClient } from "../config/database";
import { HTTPException } from "hono/http-exception";
import { UserValidation } from "../validation/user-validation";

export class UserService {
	static async get(token: User): Promise<UserResponse> {
		const user = await prismaClient.user.findFirst({
			where: {
				username: token.username,
			},
		});

		if (!user) {
			throw new HTTPException(401, {
				message: "Unauthorized",
			});
		}
		return toUserResponse(user);
	}

	static async update(
		user: User,
		request: UpdateUserRequest
	): Promise<UserResponse> {
		request = UserValidation.UPDATE.parse(request);

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

		return toUserResponse(user, true);
	}
}

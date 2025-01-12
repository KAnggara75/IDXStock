import {
	type LoginUserRequest,
	type RegisterUserRequest,
	toUserResponse,
	type UserResponse,
} from "../model/user-model.ts";
import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { prismaClient } from "../config/database.ts";
import { AuthValidation } from "../validation/auth-validation.ts";

export class AuthService {
	static async register(request: RegisterUserRequest): Promise<UserResponse> {
		request = AuthValidation.REGISTER.parse(request);

		const totalUserWithSameUsername: number = await prismaClient.user.count({
			where: {
				username: request.username,
			},
		});

		if (totalUserWithSameUsername != 0) {
			throw new HTTPException(400, {
				message: "Username already exists",
			});
		}

		request.password = await Bun.password.hash(request.password, {
			algorithm: "bcrypt",
			cost: 10,
		});

		const user = await prismaClient.user.create({
			data: request,
		});

		return toUserResponse(user);
	}

	static async login(request: LoginUserRequest): Promise<UserResponse> {
		request = AuthValidation.LOGIN.parse(request);

		const user: User | null = await prismaClient.user.findUnique({
			where: {
				username: request.username,
			},
		});

		if (!user) {
			throw new HTTPException(401, {
				message: "Username or Password is wrong",
			});
		}

		const isPasswordValid = await Bun.password.verify(
			request.password,
			user.password,
			"bcrypt"
		);
		if (!isPasswordValid) {
			throw new HTTPException(401, {
				message: "Username or password is wrong",
			});
		}
		return toUserResponse(user);
	}
}

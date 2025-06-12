import {
	type LoginUserRequest,
	type RegisterUserRequest,
	toUserResponse,
	type UserResponse,
} from "../model/user-model.ts";
import { log } from "../config/logger.ts";
import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { prismaClient } from "../config/database.ts";
import { type CustomError, toErrorDetail } from "../model/errors-model.ts";

export class AuthService {
	static async register(request: RegisterUserRequest): Promise<UserResponse> {
		const totalUserWithSameUsername: number = await prismaClient.user.count({
			where: {
				OR: [
					{
						email: request.email,
					},
					{
						username: request.username,
					},
				],
			},
		});

		if (totalUserWithSameUsername != 0) {
			const errorPayload: CustomError[] = await toErrorDetail(
				"bad_request",
				"Username or email already exists",
				["username", "email"]
			);
			throw new HTTPException(400, {
				message: JSON.stringify(errorPayload),
			});
		}

		request.password = await Bun.password.hash(request.password, {
			algorithm: "bcrypt",
			cost: 10,
		});

		const user = await prismaClient.user.create({
			data: request,
		});

		return toUserResponse(user, true);
	}

	static async login(request: LoginUserRequest): Promise<UserResponse> {
		log.debug(`${request.username} Try to login`);
		const errorPayload: CustomError[] = await toErrorDetail(
			"unauthorized",
			"Username or password is wrong",
			["username", "password"]
		);
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
		log.debug(
			`${request.username} Login result: ${user ? "success" : "failed"}`
		);

		const isPasswordValid = await Bun.password.verify(
			request.password,
			user.password,
			"bcrypt"
		);

		if (!isPasswordValid) {
			throw new HTTPException(401, {
				message: JSON.stringify(errorPayload),
			});
		}
		return toUserResponse(user, true);
	}

	static async logout(user: User): Promise<void> {
		try {
			log.debug(`${user.username} Try to logout`);
			await prismaClient.user.update({
				where: {
					username: user.username,
				},
				data: {
					logoutAt: Math.floor(Date.now() / 1000),
				},
			});
		} catch {
			log.error(`${user.username} Logout failed, user not registered`);
			throw new HTTPException(404, {
				message: "user not register",
			});
		}
	}
}

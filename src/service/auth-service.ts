import {
	type LoginUserRequest,
	type RegisterUserRequest,
	toUserResponse,
	type UserResponse,
} from "../model/user-model";
import { log } from "../config/logger";
import { prismaClient } from "../config/database";
import { RedisService } from "./redis-service";
import type { Prisma, User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { type CustomError, toErrorDetail } from "../model/errors-model";

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

		const user: User = await prismaClient.user.create({
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

		const isPasswordValid: boolean = await Bun.password.verify(
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
		const now: number = Math.floor(Date.now() / 1000);
		log.debug(`${user.username} Try to logout, set logoutAt=${now}`);
		const result: Prisma.BatchPayload = await prismaClient.user.updateMany({
			where: {
				username: user.username,
			},
			data: {
				logoutAt: now,
			},
		});

		if (result.count !== 1) {
			log.warn(`${user.username} Logout failed, user not registered`);
			throw new HTTPException(403, {
				message: `Logout forbidden: user ${user.username} not registered or already revoked`,
			});
		} else {
			await RedisService.setLogoutAt(user.username, now);
		}
	}
}

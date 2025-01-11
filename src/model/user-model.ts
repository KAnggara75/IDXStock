import type { User } from "@prisma/client";
import { JwtHelper } from "../helpers/jwt-helper";

export type RegisterUserRequest = {
	email: string;
	username: string;
	password: string;
	name: string;
};

export type LoginUserRequest = {
	username: string;
	password: string;
};

export type UpdateUserRequest = {
	password?: string;
	name?: string;
};

export type UserResponse = {
	email: string;
	username: string;
	name: string;
	token?: string;
};

export async function toUserResponse(user: User): Promise<UserResponse> {
	const token = await JwtHelper.jwtSign(user);
	return {
		email: user.email,
		name: user.name,
		username: user.username,
		token: token,
	};
}

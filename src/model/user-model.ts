import type { User } from "@prisma/client";

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

export function toUserResponse(user: User): UserResponse {
	return {
		email: user.email,
		name: user.name,
		username: user.username,
	};
}

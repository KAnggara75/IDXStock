import type { User } from "@prisma/client";
import { sign } from "hono/jwt";

export class JwtHelper {
	static async jwtSign(user: User): Promise<string> {
		const secret = Bun.env.JWT_SECERET_KEY ?? "";

		const payload = {
			email: user.email,
			name: user.name,
			username: user.username,
			exp: Math.floor(Date.now() / 1000) + 60 * 5, // Token expires in 5 minutes
		};
		const token = await sign(payload, secret);
		return token;
	}
}

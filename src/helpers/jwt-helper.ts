import type { User } from "@prisma/client";
import { sign, verify } from "hono/jwt";

export class JwtHelper {
	static async jwtSign(user: User): Promise<string> {
		const secret = Bun.env.JWT_SECERET_KEY ?? "";

		const payload = {
			email: user.email,
			name: user.name,
			username: user.username,
			exp: Math.floor(Date.now() / 1000) + 60 * 60 * 24 * 30, // Token expires in 30 days
		};
		const token = await sign(payload, secret);
		return token;
	}

	static async jwtVerivy(token: string) {
		const secretKey = Bun.env.JWT_SECERET_KEY ?? "";

		const decodedPayload = await verify(token, secretKey);
		// eslint-disable-next-line no-console
		console.log(decodedPayload);
	}
}

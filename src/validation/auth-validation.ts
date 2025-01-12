import { z, ZodType } from "zod";

export class AuthValidation {
	static readonly REGISTER: ZodType = z.object({
		username: z.string().min(1).max(100),
		password: z.string().min(1).max(100),
		email: z.string().email(),
		name: z.string().min(1).max(100),
	});

	static readonly LOGIN: ZodType = z.object({
		username: z.string().min(1).max(100),
		password: z.string().min(1).max(100),
	});

	static readonly TOKEN: ZodType = z.string().min(1);
}

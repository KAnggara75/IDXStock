import { type Context, Hono } from "hono";
import { log } from "../config/logger";
import { AuthService } from "../service/auth-service";
import type {
	LoginUserRequest,
	RegisterUserRequest,
	UserResponse,
} from "../model/user-model";
import type { ApplicationVariables } from "../model/app-model";
import { AuthValidation } from "../validation/auth-validation";
import { validateWithSchema } from "../validation/validate-with-schema";

export const authController = new Hono<{ Variables: ApplicationVariables }>();

authController.post(
	"/register",
	validateWithSchema(AuthValidation.REGISTER),
	async (c: Context<{ Variables: ApplicationVariables }>) => {
		log.debug("Try to registering new user");
		const request: RegisterUserRequest = await c.req.json();
		const response: UserResponse = await AuthService.register(request);
		log.debug(`Registering ${response.username}`);

		return c.json({
			data: response,
		});
	}
);

authController.post(
	"/login",
	validateWithSchema(AuthValidation.LOGIN),
	async (c: Context<{ Variables: ApplicationVariables }>) => {
		const request: LoginUserRequest = await c.req.json();
		const response: UserResponse = await AuthService.login(request);

		return c.json({
			data: response,
		});
	}
);

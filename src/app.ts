import { Hono } from "hono";
import pkg from "../package.json";
import { log } from "./config/logger";
import { requestLogger } from "./middleware/request-middleware";
import { jsonMiddleware, validateJsonBody } from "./middleware/json-middleware";
import { authController } from "./controller/auth-controller";
import { authMiddleware } from "./middleware/auth-middleware";
import { userController } from "./controller/user-controller";
import { stockController } from "./controller/stock-controller";
import { converterController } from "./controller/convert-controller";
import { errorHandler } from "./middleware/error-handler";
// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-ignore
import type { BlankEnv, BlankSchema } from "hono/dist/types/types";

export const app: Hono<BlankEnv, BlankSchema, "/api"> = new Hono().basePath(
	"/api"
);

app.use(requestLogger);
app.use(validateJsonBody);

app.get("/version", (c) => {
	log.debug("get /version");
	return c.json(
		{
			message: pkg.description ?? "KAnggara Web APP",
			version: Bun.env.APP_VERSION ?? pkg.version ?? "0.0.1",
			stability: Bun.env.APP_STAB ?? "Developer-Preview",
		},
		200
	);
});

app.route("/auth", authController);

app.use(jsonMiddleware);
app.use(authMiddleware);

app.route("/users", userController);
app.route("/stocks", stockController);
app.route("/convert", converterController);

app.notFound((c) => {
	log.warn(`Not Found: ${c.req.path}`);
	return c.json({ errors: "Not Found" }, 404);
});

app.onError(errorHandler);

export default app;

import { Hono } from "hono";
import type { ApplicationVariables } from "../model/app-model";
import { ConvertUsecase } from "../usecase/converter-usecase.ts";
import type { User } from "@prisma/client";

export const converterController = new Hono<{
	Variables: ApplicationVariables;
}>();

converterController.post("/convert", async (c) => {
	const body = await c.req.parseBody();

	const xlxsFile: string | File = body["file"];
	const user: User = c.get("user");

	const [response, filename] = await ConvertUsecase.processExcelSummary(
		xlxsFile,
		user
	);

	return c.json({
		data: {
			csv: `https://${new URL(c.req.url).host}/${filename}`,
			data: response,
		},
	});
});

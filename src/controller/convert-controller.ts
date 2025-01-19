import { Hono } from "hono";
import type { ApplicationVariables } from "../model/app-model";
import { ConvertUsecase } from "../usecase/converter-usecase.ts";
import type { SummaryData } from "../model/summary-model.ts";

export const converterController = new Hono<{
	Variables: ApplicationVariables;
}>();

converterController.post("/convert", async (c) => {
	const body = await c.req.parseBody();

	const xlxsFile: string | File = body["file"];

	const response: SummaryData[] =
		await ConvertUsecase.processExcelSummary(xlxsFile);

	return c.json({
		data: response,
	});
});

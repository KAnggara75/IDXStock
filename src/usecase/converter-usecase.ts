import { FileValidation } from "../validation/file-validation";
import { SheetService } from "../service/spreadsheet-service";
import { HTTPException } from "hono/http-exception";
import type { StockDB, SummaryData } from "../model/summary-model";
import type { User } from "@prisma/client";
import { CsvService } from "../service/csv-service.ts";

export class ConvertUsecase {
	static async processExcelSummary(
		request: string | File,
		user: User
	): Promise<(StockDB[] | string)[]> {
		// validate is user send correct file type
		const xlxsFile: File = FileValidation.DOCUMENT_CHECK.parse(request);

		try {
			// Convert .xlxs to CSV
			const fileName: string = await SheetService.toCsv(
				user.username,
				xlxsFile
			);

			// Convert CSV to Json
			const jsonValue: SummaryData[] = await CsvService.toJson(fileName);
			const response: StockDB[] = await CsvService.toJsonModel(jsonValue);

			return [response, fileName];
		} catch {
			throw new HTTPException(400, {
				message: "bad request file name",
			});
		}
	}
}

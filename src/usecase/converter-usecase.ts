import { FileValidation } from "../validation/file-validation";
import { SheetService } from "../service/spreadsheet-service";
import { HTTPException } from "hono/http-exception";
import type { SummaryData } from "../model/summary-model";

export class ConvertUsecase {
	static async processExcelSummary(
		request: string | File
	): Promise<SummaryData[]> {
		// validate is user send correct file type
		const xlxsFile: File = FileValidation.DOCUMENT_CHECK.parse(request);

		try {
			// Convert .xlxs to SummaryData[]
			return await SheetService.toSummaryData(xlxsFile);
		} catch {
			throw new HTTPException(400, {
				message: "bad request file name",
			});
		}
	}
}

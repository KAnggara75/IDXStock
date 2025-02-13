import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { DailyRepository } from "../repository/daily-repo";
import { GoogleFinanceService } from "../service/googleFinance-service";
import { type StockModel, toGoogleFinance } from "../model/googleFinance-model";
import type { SummaryData } from "../model/summary-model.ts";
import { FileValidation } from "../validation/file-validation.ts";
import { SheetService } from "../service/spreadsheet-service.ts";
import { HistoryRepository } from "../repository/history-repository.ts";
import { log } from "../config/logger.ts";

export class StockUsecase {
	static async updateDailyStock(user: User): Promise<StockModel[]> {
		try {
			const stockResponse: Response =
				await GoogleFinanceService.getStocksSummary();
			let stockData: StockModel[] = [];

			if (stockResponse.ok) {
				stockData = toGoogleFinance(await stockResponse.json());
				await DailyRepository.upsert(stockData, user);
			}
			return stockData;
		} catch {
			throw new HTTPException(400, {
				message: "bad request file name",
			});
		}
	}

	static async addStockSummary(
		user: User,
		request: string | File
	): Promise<SummaryData[]> {
		// validate is user send correct file type
		const xlxsFile: File = FileValidation.DOCUMENT_CHECK.parse(request);

		// Convert .xlxs to SummaryData[]
		const data: SummaryData[] = await SheetService.toSummaryData(xlxsFile);

		HistoryRepository.insertSummary(data).catch((e) => log.error(e));

		return data;
	}
}

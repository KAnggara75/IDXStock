import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { DailyRepository } from "../repository/daily-repository";
import { GoogleFinanceService } from "../service/googleFinance-service";
import { type StockModel, toGoogleFinance } from "../model/googleFinance-model";
import type { SummaryData } from "../model/summary-model";
import { FileValidation } from "../validation/file-validation";
import { SheetService } from "../service/spreadsheet-service";
import { HistoryRepository } from "../repository/history-repository";
import { log } from "../config/logger";

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

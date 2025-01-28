import type { User } from "@prisma/client";
import { HTTPException } from "hono/http-exception";
import { DailyRepository } from "../repository/daily-repo";
import { GoogleFinanceService } from "../service/googleFinance-service";
import { type StockModel, toGoogleFinance } from "../model/googleFinance-model";
import type { SummaryData } from "../model/summary-model.ts";
import { FileValidation } from "../validation/file-validation.ts";
import { SheetService } from "../service/spreadsheet-service.ts";
import { HistoryRepository } from "../repository/history-repository.ts";

export class StockUsecase {
	static async updateDailyStock(user: User): Promise<StockModel[]> {
		try {
			const stockResponse: Response =
				await GoogleFinanceService.getStocksSummary();
			let stockData: StockModel[] = [];

			if (stockResponse.ok) {
				stockData = toGoogleFinance(await stockResponse.json());
				await DailyRepository.upsert(stockData, user);
				return stockData;
			} else {
				throw new HTTPException(400, {
					message: "bad request file name",
				});
			}
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

		try {
			// Convert .xlxs to SummaryData[]
			const data: SummaryData[] = await SheetService.toSummaryData(xlxsFile);

			HistoryRepository.insertSummary(data);

			return data;
		} catch {
			throw new HTTPException(400, {
				message: "bad request file name",
			});
		}
	}
}

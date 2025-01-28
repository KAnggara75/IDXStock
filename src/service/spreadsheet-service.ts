import type { WorkBook, WorkSheet } from "xlsx";
import * as XLSX from "xlsx";
import type { SummaryData } from "../model/summary-model.ts";

export class SheetService {
	static async toSummaryData(request: File): Promise<SummaryData[]> {
		return await request.arrayBuffer().then((res) => {
			const data = new Uint8Array(res);
			const workbook: WorkBook = XLSX.read(data, { type: "array" });

			const worksheet: WorkSheet = workbook.Sheets[workbook.SheetNames[0]];

			const data_headers: string[] = [
				"no",
				"stock_code",
				"company_name",
				"remarks",
				"previous",
				"open_price",
				"last_trading_date",
				"first_trade",
				"high",
				"low",
				"close",
				"change",
				"volume",
				"value",
				"frequency",
				"index_individual",
				"offer",
				"offer_volume",
				"bid",
				"bid_volume",
				"listed_shares",
				"tradeble_shares",
				"weight_for_index",
				"foreign_sell",
				"foreign_buy",
				"non_regular_volume",
				"non_regular_value",
				"non_regular_frequency",
			];

			const jsonData: SummaryData[] = XLSX.utils.sheet_to_json(worksheet, {
				header: data_headers,
				raw: false,
				defval: null,
			});

			return jsonData.slice(1);
		});
	}
}

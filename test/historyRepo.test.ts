import { beforeEach, describe, it } from "bun:test";
import { HistoryRepository } from "../src/repository/history-repository.ts";
import type { SummaryData } from "../src/model/summary-model.ts";
import { HistoryTest } from "./test-util.ts";

describe("HistoryRepo", () => {
	beforeEach(async () => {
		await HistoryTest.delete();
	});

	it("Test", async () => {
		const summaryData: SummaryData = {
			bid: 0,
			bid_volume: 0,
			change: 0,
			close: 0,
			company_name: "Unit TEST Company Name",
			first_trade: 0,
			foreign_buy: 0,
			foreign_sell: 0,
			frequency: 0,
			high: 0,
			index_individual: "1.2",
			last_trading_date: "24 Jan 1999",
			listed_shares: 0,
			low: 0,
			non_regular_frequency: 0,
			non_regular_value: 0,
			non_regular_volume: 0,
			offer: 0,
			offer_volume: 0,
			open_price: 0,
			previous: 0,
			remarks: "",
			stock_code: "TEST",
			tradeble_shares: 0,
			value: 0,
			volume: 0,
			weight_for_index: 0,
		};

		const data: SummaryData[] = [summaryData];

		await HistoryRepository.insertSummary(data);
	});
});

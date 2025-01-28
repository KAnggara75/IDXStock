import { describe, expect, it } from "bun:test";
import { SheetService } from "../src/service/spreadsheet-service.ts";
import type { SummaryData } from "../src/model/summary-model.ts";

describe("Convert xlsx to JSON", () => {
	it("Code %s, price %d ", async () => {
		const path = "./test/Stock Summary-12345678.xlsx";
		const blob = Bun.file(path);
		const file = new File([blob], "Stock Summary-12345678.xlsx", {
			type: blob.type,
		});

		const data: SummaryData[] = await SheetService.toSummaryData(file);

		expect(data[0].stock_code).toBe("AADI");
		expect(data[1].stock_code).toBe("AALI");
		expect(data[2].stock_code).toBe("ZYRX");

		expect(Number(data[0].previous)).toBe(8800);
		expect(Number(data[1].previous)).toBe(5975);
		expect(Number(data[2].previous)).toBe(131);

		expect(Number(data[0].listed_shares)).toBe(7786891760);
		expect(Number(data[1].listed_shares)).toBe(1924688333);
		expect(Number(data[2].listed_shares)).toBe(1333334556);

		expect(Number(data[0].close)).toBe(8900);
		expect(Number(data[1].close)).toBe(6000);
		expect(Number(data[2].close)).toBe(132);

		expect(Number(data[0].high)).toBe(9200);
		expect(Number(data[1].high)).toBe(6000);
		expect(Number(data[2].high)).toBe(132);

		expect(Number(data[0].low)).toBe(8825);
		expect(Number(data[1].low)).toBe(5900);
		expect(Number(data[2].low)).toBe(130);
	});
});

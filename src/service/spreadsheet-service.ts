import * as XLSX from "xlsx";

export class SheetService {
	static async toCsv(username: string, request: File): Promise<string> {
		return await request.arrayBuffer().then((res) => {
			const data = new Uint8Array(res);
			const workbook = XLSX.read(data, { type: "array" });

			const first_sheet_name = workbook.SheetNames[0];

			const worksheet = workbook.Sheets[first_sheet_name];

			const jsonData = XLSX.utils.sheet_to_json(worksheet, {
				raw: false,
				defval: null,
			});

			const outputCsv = `public/csv/${username}-${request.name.substring(0, request.name.indexOf(".")).split("-")[1]}.csv`;

			const new_worksheet = XLSX.utils.json_to_sheet(jsonData);
			const new_workbook = XLSX.utils.book_new();

			XLSX.utils.book_append_sheet(new_workbook, new_worksheet, "csv_sheet");

			XLSX.writeFile(new_workbook, outputCsv);

			return outputCsv;
		});
	}
}

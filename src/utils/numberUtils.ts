export class NumberUtil {
	static toInt(number: string): number | null {
		return number === "" ? null : parseInt(number);
	}

	static toFloat(number: string): number | null {
		return number === "" ? null : parseFloat(number.replaceAll(",", "."));
	}
}

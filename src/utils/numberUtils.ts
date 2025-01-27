export class NumberUtil {
	static toInt(number: string | number): number | null {
		if (typeof number === "number") {
			return number;
		}

		return number === "" ? null : parseInt(number);
	}

	static toFloat(number: string | number): number | null {
		if (typeof number === "number") {
			return number;
		}

		return number === "" ? null : parseFloat(number.replaceAll(",", "."));
	}

	static toBigInt(number: string | number): bigint | null {
		if (typeof number === "number") {
			return BigInt(number);
		}

		if (number.includes("E+")) {
			return BigInt(Number(number));
		}

		return number === "" ? null : BigInt(number);
	}
}

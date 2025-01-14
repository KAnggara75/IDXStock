export class DateUtils {
	/**
	 * Returns ISO-8691 Date Format from
	 * @param dateString dateTime with format YYYY-MM-DD hh:mm:ss
	 */
	static toDate(dateString: string): Date {
		try {
			if (dateString.length != 19) {
				return new Date();
			}

			return new Date(
				new Date(dateString.replace(" ", "T")).getTime() + 7 * 60 * 60 * 1000
			);
		} catch {
			return new Date();
		}
	}

	/**
	 * Returns Token Expired in second
	 * @param days from now to token Expired.
	 */
	static expInDays(days: number | string): number {
		const now = Math.floor(Date.now() / 1000);
		let day: number = 1;

		if (!isNaN(Number(days))) {
			day = Number(days);
		}

		const future: number = 3600 * 24 * day;

		return now + future;
	}
}

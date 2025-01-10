export class DateUtils {
	static toDate(dateString: string): Date {
		try {
			const date = new Date(dateString.replace(" ", "T"));
			const offset = 7; // Waktu Jakarta adalah UTC+7
			const localDate = new Date(date.getTime() + offset * 60 * 60 * 1000);
			return localDate;
		} catch {
			return new Date();
		}
	}
}

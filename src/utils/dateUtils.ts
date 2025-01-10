export class DateUtils {
	static toDate(dateString: string): Date {
		try {
			const date = new Date(dateString.replace(" ", "T"));
			return date;
		} catch {
			return new Date();
		}
	}
}

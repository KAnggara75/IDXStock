export class JsonUtils {
	static safeParseJSON(input: string) {
		try {
			return JSON.parse(input);
		} catch {
			return input;
		}
	}
}

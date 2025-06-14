const envKeywords = Bun.env.SENSITIVE_SQL_KEYWORDS || "";
const sensitiveKeywords: string[] = envKeywords
	.split(",")
	.map((k) => k.trim().toLowerCase())
	.filter(Boolean);

export function isSensitiveQuery(query: string): boolean {
	const lowerQuery = query.toLowerCase();
	return sensitiveKeywords.some((k) => lowerQuery.includes(k));
}

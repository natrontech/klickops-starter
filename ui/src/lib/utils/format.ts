export function formatBytes(bytes: number): string {
	if (bytes < 1024) return `${bytes} B`;
	const units = ["KiB", "MiB", "GiB"];
	let value = bytes;
	let unit = "B";
	for (const next of units) {
		if (value < 1024) break;
		value /= 1024;
		unit = next;
	}
	return `${value.toFixed(value >= 10 ? 0 : 1)} ${unit}`;
}

export function formatDate(iso: string): string {
	return new Date(iso).toLocaleString(undefined, {
		dateStyle: "medium",
		timeStyle: "short"
	});
}

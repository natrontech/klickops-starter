import { describe, expect, it } from "vitest";
import { formatBytes } from "./format";

describe("formatBytes", () => {
	it("keeps small values in bytes", () => {
		expect(formatBytes(0)).toBe("0 B");
		expect(formatBytes(512)).toBe("512 B");
	});

	it("converts to binary units", () => {
		expect(formatBytes(1024)).toBe("1.0 KiB");
		expect(formatBytes(5 * 1024 * 1024)).toBe("5.0 MiB");
		expect(formatBytes(1536)).toBe("1.5 KiB");
	});

	it("drops decimals for large values", () => {
		expect(formatBytes(25 * 1024 * 1024)).toBe("25 MiB");
	});
});

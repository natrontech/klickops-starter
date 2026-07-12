// Single fetch wrapper for the backend API. Every API module goes
// through this - never call fetch("/api/...") directly in a component.
export class ApiError extends Error {
	constructor(
		message: string,
		public status: number
	) {
		super(message);
	}
}

export async function api<T>(path: string, init?: RequestInit): Promise<T> {
	const res = await fetch(`/api${path}`, init);
	if (!res.ok) {
		let message = res.statusText;
		try {
			const body = await res.json();
			if (typeof body.error === "string") message = body.error;
		} catch {
			// non-JSON error body - keep the status text
		}
		throw new ApiError(message, res.status);
	}
	if (res.status === 204) return undefined as T;
	return res.json();
}

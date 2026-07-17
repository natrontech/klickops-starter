import { api } from "./client";

export interface Visits {
	visits: number;
}

export function fetchVisits(): Promise<Visits> {
	return api<Visits>("/visits");
}

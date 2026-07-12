import { api } from "./client";

export interface Health {
	status: string;
	database: boolean;
	storage: boolean;
}

export function fetchHealth(): Promise<Health> {
	return api<Health>("/healthz");
}

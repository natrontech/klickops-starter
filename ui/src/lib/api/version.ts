import { api } from "./client";

export interface Build {
	env: string;
	gitSha: string;
}

export function fetchVersion(): Promise<Build> {
	return api<Build>("/version");
}

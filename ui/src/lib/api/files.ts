import { api } from "./client";

export interface FileInfo {
	key: string;
	size: number;
	lastModified: string;
}

export function listFiles(): Promise<FileInfo[]> {
	return api<FileInfo[]>("/files");
}

export function uploadFile(file: File): Promise<FileInfo> {
	const form = new FormData();
	form.append("file", file);
	return api<FileInfo>("/files", { method: "POST", body: form });
}

export function deleteFile(key: string): Promise<void> {
	return api<void>(`/files/${encodeURIComponent(key)}`, { method: "DELETE" });
}

export function fileUrl(key: string): string {
	return `/api/files/${encodeURIComponent(key)}`;
}

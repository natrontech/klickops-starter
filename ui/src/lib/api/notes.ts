import { api } from "./client";

export interface Note {
	id: number;
	text: string;
	createdAt: string;
}

export function listNotes(): Promise<Note[]> {
	return api<Note[]>("/notes");
}

export function createNote(text: string): Promise<Note> {
	return api<Note>("/notes", {
		method: "POST",
		headers: { "Content-Type": "application/json" },
		body: JSON.stringify({ text })
	});
}

export function deleteNote(id: number): Promise<void> {
	return api<void>(`/notes/${id}`, { method: "DELETE" });
}

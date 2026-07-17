<script lang="ts">
	import Badge from "$lib/components/ui/badge.svelte";
	import Button from "$lib/components/ui/button.svelte";
	import Card from "$lib/components/ui/card.svelte";
	import Input from "$lib/components/ui/input.svelte";
	import { fetchHealth, type Health } from "$lib/api/health";
	import { createNote, deleteNote, listNotes, type Note } from "$lib/api/notes";
	import { deleteFile, fileUrl, listFiles, uploadFile, type FileInfo } from "$lib/api/files";
	import { fetchVisits } from "$lib/api/visits";
	import { formatBytes, formatDate } from "$lib/utils/format";

	let health = $state<Health | null>(null);
	let notes = $state<Note[]>([]);
	let files = $state<FileInfo[]>([]);
	let visits = $state<number | null>(null);
	let newNote = $state("");
	let noteError = $state("");
	let fileError = $state("");
	let uploading = $state(false);

	$effect(() => {
		void load();
	});

	async function load() {
		try {
			health = await fetchHealth();
		} catch {
			health = null;
			return;
		}
		if (health.database) notes = await listNotes().catch(() => []);
		if (health.storage) files = await listFiles().catch(() => []);
		if (health.cache) visits = (await fetchVisits().catch(() => null))?.visits ?? null;
	}

	async function addNote(event: SubmitEvent) {
		event.preventDefault();
		noteError = "";
		try {
			const note = await createNote(newNote);
			notes = [note, ...notes];
			newNote = "";
		} catch (err) {
			noteError = err instanceof Error ? err.message : "failed to add note";
		}
	}

	async function removeNote(id: number) {
		await deleteNote(id).catch(() => {});
		notes = notes.filter((n) => n.id !== id);
	}

	async function onFilePicked(event: Event) {
		const input = event.currentTarget as HTMLInputElement;
		const file = input.files?.[0];
		if (!file) return;
		fileError = "";
		uploading = true;
		try {
			await uploadFile(file);
			files = await listFiles();
		} catch (err) {
			fileError = err instanceof Error ? err.message : "upload failed";
		} finally {
			uploading = false;
			input.value = "";
		}
	}

	async function removeFile(key: string) {
		await deleteFile(key).catch(() => {});
		files = files.filter((f) => f.key !== key);
	}
</script>

<div class="flex flex-col gap-6">
	<div>
		<h1 class="text-[28px] font-semibold leading-tight tracking-[-0.015em]">It runs.</h1>
		<p class="mt-2 max-w-lg text-sm text-muted-foreground">
			This page is served by the Go backend from a single container. Below it talks to PostgreSQL
			and S3 storage - bind them as services on klickops, then replace all of this with your own
			app.
		</p>
	</div>

	<div class="flex flex-wrap gap-2">
		{#if health}
			<Badge tone="success">API connected</Badge>
			<Badge tone={health.database ? "success" : "muted"}>
				Database {health.database ? "connected" : "not bound"}
			</Badge>
			<Badge tone={health.storage ? "success" : "muted"}>
				Storage {health.storage ? "connected" : "not bound"}
			</Badge>
			<Badge tone={health.cache ? "success" : "muted"}>
				Cache {health.cache ? "connected" : "not bound"}
			</Badge>
		{:else}
			<Badge tone="warning">API unreachable - is the backend running? (make dev-backend)</Badge>
		{/if}
	</div>

	<Card title="Notes" description="A tiny CRUD example backed by PostgreSQL.">
		{#if health?.database}
			<form onsubmit={addNote} class="flex gap-2">
				<Input bind:value={newNote} placeholder="Write a note…" maxlength={2000} />
				<Button type="submit" disabled={!newNote.trim()}>Add</Button>
			</form>
			{#if noteError}
				<p class="mt-2 text-xs text-destructive">{noteError}</p>
			{/if}
			<ul class="mt-4 flex flex-col divide-y divide-border/60">
				{#each notes as note (note.id)}
					<li class="flex items-center justify-between gap-3 py-2.5">
						<div class="min-w-0">
							<p class="truncate text-sm">{note.text}</p>
							<p class="text-[11px] text-muted-foreground">{formatDate(note.createdAt)}</p>
						</div>
						<Button variant="ghost" onclick={() => removeNote(note.id)}>Delete</Button>
					</li>
				{:else}
					<li class="py-2.5 text-sm text-muted-foreground">No notes yet - add one above.</li>
				{/each}
			</ul>
		{:else}
			<p class="text-sm text-muted-foreground">
				No database bound. On klickops: add a <strong>PostgreSQL</strong> service to your project
				and bind it to this app as <code class="font-mono text-xs">DATABASE_URL</code>. Locally:
				<code class="font-mono text-xs">docker compose up -d</code> and copy
				<code class="font-mono text-xs">.env.example</code> to
				<code class="font-mono text-xs">.env</code>.
			</p>
		{/if}
	</Card>

	<Card title="Files" description="Upload and download files from S3-compatible storage.">
		{#if health?.storage}
			<label
				class="inline-flex h-8 cursor-pointer items-center rounded-lg bg-primary px-3 text-xs font-medium text-primary-foreground transition-colors hover:bg-primary/90 {uploading
					? 'pointer-events-none opacity-50'
					: ''}"
			>
				{uploading ? "Uploading…" : "Upload a file"}
				<input type="file" class="hidden" onchange={onFilePicked} disabled={uploading} />
			</label>
			{#if fileError}
				<p class="mt-2 text-xs text-destructive">{fileError}</p>
			{/if}
			<ul class="mt-4 flex flex-col divide-y divide-border/60">
				{#each files as file (file.key)}
					<li class="flex items-center justify-between gap-3 py-2.5">
						<div class="min-w-0">
							<a
								href={fileUrl(file.key)}
								target="_blank"
								rel="noreferrer"
								class="truncate text-sm hover:underline">{file.key}</a
							>
							<p class="text-[11px] text-muted-foreground">
								{formatBytes(file.size)} · {formatDate(file.lastModified)}
							</p>
						</div>
						<Button variant="ghost" onclick={() => removeFile(file.key)}>Delete</Button>
					</li>
				{:else}
					<li class="py-2.5 text-sm text-muted-foreground">No files yet - upload one above.</li>
				{/each}
			</ul>
		{:else}
			<p class="text-sm text-muted-foreground">
				No storage bound. On klickops: add a <strong>Bucket</strong> service to your project and
				bind its endpoint and credentials as
				<code class="font-mono text-xs">S3_*</code> variables (see README). Locally:
				<code class="font-mono text-xs">docker compose up -d</code> starts an S3-compatible server.
			</p>
		{/if}
	</Card>

	<Card
		title="Cache"
		description="A Valkey (Redis-compatible) example: a visit counter, plus the notes list above is served cache-aside with a 30s TTL."
	>
		{#if health?.cache}
			<p class="text-sm">
				You are visitor
				<span class="font-mono font-semibold">#{visits ?? "…"}</span> - counted with an atomic
				<code class="font-mono text-xs">INCR</code>. Reload to bump it. The notes list sets an
				<code class="font-mono text-xs">X-Cache: hit|miss</code> response header - watch it flip
				in the network tab.
			</p>
		{:else}
			<p class="text-sm text-muted-foreground">
				No cache bound. On klickops: add a <strong>Valkey</strong> service to your project and
				connect it to this app - it injects
				<code class="font-mono text-xs">REDIS_URL</code>. Locally:
				<code class="font-mono text-xs">docker compose up -d</code> starts a Valkey server.
			</p>
		{/if}
	</Card>

	<Card title="Next steps">
		<ol class="list-inside list-decimal space-y-1.5 text-sm text-muted-foreground">
			<li>Open this folder with Claude Code (or Cursor, or any AI coding tool).</li>
			<li>Describe the app you want - the AI knows this codebase via CLAUDE.md.</li>
			<li>Push to GitHub and deploy on klickops: connect the repo, it builds and ships.</li>
		</ol>
	</Card>
</div>

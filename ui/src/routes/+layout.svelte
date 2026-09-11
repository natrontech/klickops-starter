<script lang="ts">
	import "../app.css";
	import type { Snippet } from "svelte";
	import Button from "$lib/components/ui/button.svelte";
	import { fetchVersion, type Build } from "$lib/api/version";

	let { children }: { children: Snippet } = $props();

	let build = $state<Build | null>(null);

	$effect(() => {
		fetchVersion().then(
			(b) => (build = b),
			() => {}
		);
	});

	// ponytail: native <dialog> + localStorage flag, no modal library, no store
	let welcome = $state<HTMLDialogElement>();

	$effect(() => {
		if (!localStorage.getItem("welcomed")) welcome?.showModal();
	});

	function dismiss() {
		localStorage.setItem("welcomed", "1");
		welcome?.close();
	}
</script>

<dialog
	bind:this={welcome}
	onclose={dismiss}
	class="m-auto max-w-sm rounded-xl border border-border bg-card p-5 text-foreground backdrop:bg-black/50"
>
	<h2 class="text-[15px] font-semibold tracking-tight">Welcome to klickops starter</h2>
	<p class="mt-2 text-sm text-muted-foreground">
		This is the demo app - notes, files and the cache counter are examples. Describe what you want
		to build and replace them.
	</p>
	<div class="mt-4 flex justify-end">
		<Button onclick={dismiss}>Got it</Button>
	</div>
</dialog>

<div class="mx-auto flex min-h-screen max-w-3xl flex-col px-6">
	<header class="flex items-center justify-between py-8">
		<a href="/" class="flex items-center gap-2.5">
			<span
				class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-primary/90 to-primary text-sm font-bold text-primary-foreground"
			>
				k
			</span>
			<span class="text-[15px] font-semibold tracking-tight">klickops starter</span>
		</a>
		<a
			href="https://github.com/natrontech/klickops-starter"
			target="_blank"
			rel="noreferrer"
			class="text-sm text-muted-foreground transition-colors hover:text-foreground"
		>
			GitHub
		</a>
	</header>

	<main class="flex-1 pb-16">
		{@render children()}
	</main>

	<footer
		class="flex flex-wrap items-center justify-between gap-3 border-t border-border/60 py-6 text-xs text-muted-foreground"
	>
		<span>
			Built from the klickops starter - open this folder with your AI coding tool and describe what
			you want to build.
		</span>
		{#if build}
			<span class="font-mono text-[11px] whitespace-nowrap">
				{build.env} · {build.gitSha}
			</span>
		{/if}
	</footer>
</div>

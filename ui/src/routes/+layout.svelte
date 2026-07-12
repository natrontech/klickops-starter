<script lang="ts">
	import "../app.css";
	import type { Snippet } from "svelte";

	let { children }: { children: Snippet } = $props();
</script>

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

	<footer class="border-t border-border/60 py-6 text-xs text-muted-foreground">
		Built from the klickops starter - open this folder with your AI coding tool and describe what
		you want to build.
	</footer>
</div>

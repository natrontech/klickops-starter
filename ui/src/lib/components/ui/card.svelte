<script lang="ts">
	import type { Snippet } from "svelte";

	type Props = {
		title?: string;
		description?: string;
		class?: string;
		children: Snippet;
	};

	let { title, description, class: className = "", children }: Props = $props();
</script>

<section class="rounded-xl border border-border bg-card p-5 {className}">
	{#if title}
		<h2 class="text-[15px] font-semibold tracking-tight">{title}</h2>
	{/if}
	{#if description}
		<p class="mt-1 text-sm text-muted-foreground">{description}</p>
	{/if}
	<div class={title || description ? "mt-4" : ""}>
		{@render children()}
	</div>
</section>

<script lang="ts">
	import type { Snippet } from "svelte";

	type Props = {
		tone?: "success" | "warning" | "muted";
		children: Snippet;
	};

	let { tone = "muted", children }: Props = $props();

	const tones = {
		success: "text-success border-success/30 bg-success/10",
		warning: "text-warning border-warning/30 bg-warning/10",
		muted: "text-muted-foreground border-border bg-muted/40"
	};
</script>

<span
	class="inline-flex items-center gap-1.5 rounded-md border px-2 py-0.5 text-[11px] font-medium {tones[
		tone
	]}"
>
	{@render children()}
</span>

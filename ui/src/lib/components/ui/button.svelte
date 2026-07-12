<script lang="ts">
	import type { Snippet } from "svelte";
	import type { HTMLButtonAttributes } from "svelte/elements";

	type Props = HTMLButtonAttributes & {
		variant?: "primary" | "outline" | "ghost" | "destructive";
		children: Snippet;
	};

	let { variant = "primary", class: className = "", children, ...rest }: Props = $props();

	const base =
		"inline-flex h-8 items-center justify-center gap-1.5 rounded-lg px-3 text-xs font-medium transition-colors focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ring disabled:pointer-events-none disabled:opacity-50";

	const variants = {
		primary: "bg-primary text-primary-foreground hover:bg-primary/90",
		outline: "border border-border bg-transparent hover:bg-muted",
		ghost: "hover:bg-muted",
		destructive: "border border-border text-destructive hover:bg-destructive/10"
	};
</script>

<button class="{base} {variants[variant]} {className}" {...rest}>
	{@render children()}
</button>

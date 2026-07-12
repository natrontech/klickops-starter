<script lang="ts">
	import type { HTMLInputAttributes } from "svelte/elements";

	type Props = HTMLInputAttributes & { value?: string };

	let { value = $bindable(""), class: className = "", ...rest }: Props = $props();
</script>

<input
	bind:value
	class="h-8 w-full rounded-lg border border-border bg-transparent px-3 text-sm placeholder:text-muted-foreground focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-ring {className}"
	{...rest}
/>

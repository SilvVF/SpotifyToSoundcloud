<script lang="ts">
    import { onDestroy, setContext, type Snippet } from "svelte";
    import { SvelteMap } from "svelte/reactivity";
    import { EventsOn } from "wailsjs/runtime/runtime";
    import type {
        GenerationState,
        MatchProgress,
    } from "src/types/types.svelte";

    const { children }: { children: Snippet } = $props()

    let generating = new SvelteMap<string, GenerationState>();

    const unsub = EventsOn("match_progress", (data: MatchProgress) => {
        console.log("received generate event: ", data)
        generating.set(data.forId, {
            total: data.total,
            complete: data.progress,
            error: data.status === "error" ? "failed to generate" : undefined,
            status: data.status,
        });
    });

    onDestroy(() => unsub());

    setContext("generateContextName", generating);
</script>

{@render children?.()}

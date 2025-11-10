<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { useQuery } from "@sveltestack/svelte-query";
    import { onMount } from "svelte";
    import { SpotifyPlaylists } from "wailsjs/go/main/App";
    import { LogDebug } from "wailsjs/runtime/runtime";

    let list = $state<string[]>([]);
    let result = useQuery({
        queryKey: ["playlists"],
        queryFn: async () => await SpotifyPlaylists(),
        placeholderData: [],
    });
</script>

<div class="flex flex-col h-full w-full scroll-auto">
    <div>HOME</div>
    <Button onclick={$result.refetch}>Fetch</Button>
    {#if !result || $result.status === "loading"}
        <span>Loading...</span>
    {:else if $result.status === "error"}
        <span>Error: {$result.error}</span>
    {:else}
        <div class="flex flex-col">
            {#each $result.data as playlist}
                <div class="text-md text-red-50">
                    {playlist}
                </div>
            {/each}
            <div>list</div>
            {#each list as playlist}
                <div class="text-md text-red-50">
                    {playlist.name}
                </div>
            {/each}
        </div>
        <div>{$result.isFetching ? "Background Updating..." : " "}</div>
    {/if}
</div>

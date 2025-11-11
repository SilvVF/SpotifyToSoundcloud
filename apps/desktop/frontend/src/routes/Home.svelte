<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { useQuery } from "@sveltestack/svelte-query";
    import { Link } from "svelte5-router";
    import { SpotifyPlaylists } from "wailsjs/go/main/App";
    import { spotify } from "wailsjs/go/models";

    let result = useQuery({
        queryKey: ["playlists"],
        queryFn: async () =>
            (await SpotifyPlaylists()) as Array<spotify.SimplePlaylist>,
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
        <div class="grid grid-cols-3 space-y-4 space-x-4">
            {#each $result.data as playlist}
                <Link to={"/convert/" + playlist.id}>
                    <div class="flex flex-col">
                        <img
                            class="size-48"
                            src={playlist.images[0].url}
                            alt="none"
                        />
                        <div>
                            {playlist.name}
                        </div>
                    </div>
                </Link>
            {/each}
        </div>
        <div>{$result.isFetching ? "Background Updating..." : " "}</div>
    {/if}
</div>

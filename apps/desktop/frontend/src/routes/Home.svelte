<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import { useQuery } from "@sveltestack/svelte-query";
    import { Link } from "svelte5-router";
    import { SpotifyPlaylists } from "wailsjs/go/main/App";

    let result = useQuery({
        queryKey: ["playlists"],
        queryFn: async () => {
            const results = await SpotifyPlaylists();
            return results;
        },
    });
</script>

<div class="flex flex-col overflow-y-auto">
    {#if !result || $result.status === "loading"}
        <span>Loading...</span>
    {:else if $result.status === "error"}
        <span>Error: {$result.error} </span>
        <Button onclick={$result.refetch}>Retry</Button>
    {:else}
        <div
            class="grid grid-cols-3 space-y-4 space-x-4 w-full h-full justify-center my-12"
        >
            {#each $result.data as playlist}
                <svelte:boundary>
                    <Link to={"/convert/" + playlist.id}>
                        <div
                            class="flex flex-col hover:bg-card items-center px-8"
                        >
                            <img
                                class="w-full object-cover aspect-square"
                                src={playlist.imgs[0].url}
                                alt="none"
                            />
                            <span
                                class="w-full text-start pt-4 text-md font-sem"
                            >
                                {playlist.title}
                            </span>
                        </div>
                    </Link>

                    {#snippet failed(error)}
                        <span>failed to load {error}</span>
                    {/snippet}
                </svelte:boundary>
            {/each}
        </div>
        <div>{$result.isFetching ? "Loading..." : " "}</div>
    {/if}
</div>

<script lang="ts">
    import Button from "$lib/components/ui/button/button.svelte";
    import { Progress } from "$lib/components/ui/progress";
    import { Spinner } from "$lib/components/ui/spinner";
    import { useQuery } from "@sveltestack/svelte-query";
    import {
        SpotifyPlaylist,
        GetMatches,
        CreateSoundCloudPlaylist,
    } from "wailsjs/go/main/App";
    import { getContext } from "svelte";
    import { toast } from "svelte-sonner";
    import type { GenerationContext } from "src/types/types.svelte";
    import { main } from "wailsjs/go/models";
    import { Skeleton } from "$lib/components/ui/skeleton";
    import { Badge } from "$lib/components/ui/badge";
    import ScrollArea from "$lib/components/ui/scroll-area/scroll-area.svelte";
    import * as Empty from "$lib/components/ui/empty/index.js";
    import SoundcloudIcon from "$lib/icons/SoundcloudIcon.svelte";
    import { EventsOn } from "wailsjs/runtime/runtime";
    import { SvelteSet } from "svelte/reactivity";
    import { cn } from "$lib/utils";
    import Separator from "$lib/components/ui/separator/separator.svelte";
    import { navigate, useLocation } from "svelte5-router";
    import { derived } from "svelte/store";

    const { playlistId }: { playlistId: string } = $props();
    const location = useLocation();
    const res = useQuery(["playlist", playlistId], async () => {
        const res = await SpotifyPlaylist(playlistId);
        return res;
    });

    const currentAnchor = derived(location, (location) => {
        return location.hash?.substring(1, location.hash.length) ?? "";
    });

    type MatchType = {
        [id: string]: {
            selected: string;
            forTrack: main.Track | undefined;
            tracks: main.ScoredTrack[];
            minScore: number;
            maxScore: number;
        };
    };

    let matchUpdates = $state(0);
    const matches = $state<MatchType>({});
    const genContext = getContext<GenerationContext>("generateContextName");

    const item = $derived(genContext.get(playlistId));
    const progress = $derived(item ? 100 * (item.complete / item.total) : 0);

    let fetchingMatches = $derived(item ? item.status === "running" : false);
    let creatingPlaylist = $state(false);

    const expanded = new SvelteSet<string>();

    function handleConfirmSelections() {
        const ids = Object.entries(matches)
            .map(([, { selected }]) => selected)
            .filter((selected) => selected !== "");

        // title string, description string, sharing string, ids []string
        const data = $res.data;
        if (data && ids.length > 0) {
            creatingPlaylist = true;
            CreateSoundCloudPlaylist(
                data.playlist.title,
                data.playlist.description,
                "private",
                ids,
            )
                .then((created: any) => {
                    console.log(created);
                    toast.success("created playlist");
                })
                .catch((e: any) =>
                    toast.error(`failed to create playlist err ${e}`),
                )
                .finally(() => {
                    creatingPlaylist = false;
                });
        }
    }

    function toggleExpanded(id: string) {
        if (expanded.has(id)) {
            expanded.delete(id);
        } else {
            expanded.add(id);
        }
    }

    function handleMatchSelected(forId: string, id: string) {
        if (matches[forId].selected === id) {
            matches[forId].selected = "";
        } else {
            matches[forId].selected = id;
        }
    }

    let canChangeAnchor = $state(true);

    $effect(() => {
        console.log(`starting observer ${matchUpdates}`);

        const observer = new IntersectionObserver(
            (entries) => {
                entries.forEach((entry) => {
                    if (entry.isIntersecting && canChangeAnchor) {
                        navigate(`#${entry.target.id}`, { replace: false });
                    }
                });
            },
            {
                root: document.getElementById("match-container"),
                rootMargin: "0px",
                threshold: 0.5,
            },
        );

        const sections = document.querySelectorAll('[id^="match_"]');
        sections.forEach((section) => {
            observer.observe(section);
        });

        return () => {
            observer.disconnect();
        };
    });

    async function handleGenerate() {
        if (!fetchingMatches) {
            fetchingMatches = true;
            try {
                const trackScores = (tracks: main.ScoredTrack[]) => {
                    let min = 0;
                    let max = 0;
                    for (const track of tracks) {
                        min = Math.min(track.score, min);
                        max = Math.max(track.score, max);
                    }
                    return {
                        min: min,
                        max: max,
                    };
                };

                const playlistTracks = $res.data?.tracks ?? [];
                const forTracks = new Map<string, main.Track>();
                for (const track of playlistTracks) {
                    forTracks.set(track.id, track);
                }

                const unsub = EventsOn(
                    `match_result_${playlistId}`,
                    (data: main.TracksWrapper) => {
                        const { min, max } = trackScores(data.tracks);
                        matches[data.for_id] = {
                            selected: data.tracks[0].track.id,
                            minScore: min,
                            maxScore: max,
                            tracks: data.tracks,
                            forTrack: forTracks.get(data.for_id),
                        };
                        matchUpdates += 1;
                    },
                );

                const matched = await GetMatches(playlistId);

                unsub();

                for (const [trackId, match] of Object.entries(matched)) {
                    const { min, max } = trackScores(match.tracks);
                    matches[trackId] = {
                        selected: match.tracks[0].track.id,
                        minScore: min,
                        maxScore: max,
                        tracks: match.tracks,
                        forTrack: forTracks.get(trackId),
                    };
                    matchUpdates += 1;
                }
            } catch (e) {
                toast.error(`failed to get matches err: ${e}`);
            }
            fetchingMatches = false;
        }
    }
</script>

{#snippet songSkeleton()}
    <div class="flex items-center space-x-4">
        <Skeleton class="size-12" />
        <div class="space-y-2 w-full">
            <Skeleton class="h-4 w-full" />
            <Skeleton class="h-4 w-3/4" />
        </div>
    </div>
{/snippet}

{#snippet trackItem(track: main.Track | undefined)}
    <div class="flex items-center space-x-4 my-2">
        <img class="size-12" src={track?.imgs[0].url} alt={track?.title} />
        <div class="space-y-2 w-full">
            <div class="h-4 w-full font-semibold text-start">
                {track?.title}
            </div>
            <div class="text-muted h-4 w-full text-start">{track?.urn}</div>
        </div>
    </div>
{/snippet}

<div class="flex flex-row w-full h-full">
    <ScrollArea class="overflow-x-clip w-1/3 border">
        {#if $res.data}
            <div class="flex flex-col mt-4 space-y-2 px-4">
                <img
                    class="size-48 object-cover"
                    src={$res.data.playlist.imgs[0].url}
                    alt=""
                />
                <span class="font-semibold">{$res.data.playlist.title}</span>
                <span class="text-muted">{$res.data.playlist.urn}</span>
                <div class="flex flex-row items-center justify-between">
                    <Button class="w-fit" onclick={handleGenerate}>
                        {#if item?.status === "running"}
                            <Spinner />
                            Generating
                        {:else}
                            Generate
                        {/if}
                    </Button>
                    {#if item?.status === "running"}
                        <span class="text-xl font-bold"
                            >{item.complete} / {item.total}</span
                        >
                    {/if}
                </div>
                <Progress value={progress} max={100} class="w-full" />
            </div>
            {#each $res.data.tracks as track}
                <div
                    class={cn(
                        $currentAnchor === `match_${track.id}` ? "bg-card" : "",
                        "px-4 py-0.5",
                    )}
                >
                    <a
                        href={`#match_${track.id}`}
                        onclick={() => {
                            canChangeAnchor = false;
                            setTimeout(() => (canChangeAnchor = true), 100);
                        }}
                    >
                        <div class={cn("flex items-center space-x-4 my-4")}>
                            <img
                                class="size-12"
                                src={track.imgs[0].url}
                                alt={track.title}
                            />
                            <div
                                class="flex flex-row items-center justify-between w-full"
                            >
                                <div class="space-y-2 w-full">
                                    <div
                                        class="flex flex-row justify-start items-start space-x-4"
                                    >
                                        <div class="h-4 w-fit font-semibold">
                                            {track.title}
                                        </div>
                                        {#if matches[track.id]}
                                            <Badge
                                                class="w-fit"
                                                variant={matches[track.id]
                                                    .maxScore > 1.25
                                                    ? "default"
                                                    : "destructive"}
                                            >
                                                Confidence {matches[track.id]
                                                    .maxScore > 1.25
                                                    ? "high"
                                                    : "low"}
                                            </Badge>
                                            <Badge
                                                class="w-fit"
                                                variant={matches[track.id]
                                                    .selected !== ""
                                                    ? "outline"
                                                    : "destructive"}
                                            >
                                                {matches[track.id].selected !==
                                                ""
                                                    ? "selected"
                                                    : "no selection"}
                                            </Badge>
                                        {:else}
                                            <Badge class="w-fit">
                                                <Spinner />
                                                Generating
                                            </Badge>
                                        {/if}
                                    </div>
                                    <div class="text-muted h-4 w-full">
                                        {track.urn}
                                    </div>
                                </div>
                                <Button
                                    onclick={() => toggleExpanded(track.id)}
                                    variant="link">stats</Button
                                >
                            </div>
                        </div>
                        {#if expanded.has(track.id)}
                            {#if matches[track.id]}
                                <span
                                    >max score: {matches[
                                        track.id
                                    ].maxScore.toFixed(2)}</span
                                >
                                <span
                                    >min score: {matches[
                                        track.id
                                    ].minScore.toFixed(2)}</span
                                >
                            {:else}
                                <Skeleton class="h-4 w-3/4" />
                                <Skeleton class="h-4 w-3/4" />
                            {/if}
                        {/if}
                    </a>
                </div>
            {/each}
        {:else}
            <div class="flex flex-col space-y-4">
                <Skeleton class="size-48" />
                {#each [0, 0, 0, 0]}
                    {@render songSkeleton()}
                {/each}
            </div>
        {/if}
    </ScrollArea>
    <div
        id="match-container"
        class="flex flex-col w-2/3 h-full border overflow-x-clip overflow-y-auto p-2"
    >
        {#if item?.status === undefined}
            <Empty.Root>
                <Empty.Header>
                    <Empty.Media variant="icon">
                        <SoundcloudIcon />
                    </Empty.Media>
                    <Empty.Title>No Matches Yet</Empty.Title>
                    <Empty.Description>
                        You haven't generated the matches for this playlist yet.
                        Click below to run the generation job.
                    </Empty.Description>
                </Empty.Header>
                <Empty.Content>
                    <div class="flex gap-2">
                        <Button onclick={handleGenerate}
                            >Generate matches</Button
                        >
                    </div>
                </Empty.Content>
            </Empty.Root>
        {:else}
            {#each Object.entries(matches) as [id, match]}
                <div id={`match_${id}`} class="flex flex-col">
                    {@render trackItem(match.forTrack)}
                    <Separator />
                    {#each match.tracks as track}
                        <button
                            type="button"
                            class={cn(
                                match.selected === track.track.id
                                    ? "bg-card"
                                    : "",
                                "w-full",
                            )}
                            onclick={() =>
                                handleMatchSelected(id, track.track.id)}
                        >
                            <div class="ms-16">
                                {@render trackItem(track.track)}
                            </div>
                        </button>
                    {/each}
                </div>
            {/each}
        {/if}
    </div>
</div>

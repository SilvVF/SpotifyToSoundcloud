<script lang="ts">
    import { GetStreams } from "wailsjs/go/main/App";
    import type { main } from "wailsjs/go/models";
    import Button from "./components/ui/button/button.svelte";
    import { Play, X } from "@lucide/svelte";
    import Spinner from "./components/ui/spinner/spinner.svelte";

    const {
        track,
    }: {
        track: main.Track;
    } = $props();

    type LoadState = "loading" | "idle" | "error" | "done";

    let loadState = $state<LoadState>("idle");
    let audioEl = $state<HTMLAudioElement | undefined>(undefined);
    let closed = $state(false);
    const audioLevelKey = "audio_level";
    let audioLevel = Number(localStorage.getItem(audioLevelKey) ?? 0.6);

    $effect(() => {
        if (!audioEl) return;
        audioEl.addEventListener("volumechange", () => {
            localStorage.setItem(audioLevelKey, `${audioEl?.volume ?? 0.6}`);
        });
    });

    async function fetchAudioWithHeaders(
        url: string,
        headers: Record<string, string>,
    ) {
        const response = await fetch(url, {
            method: "GET",
            headers: {
                Authorization: headers["Authorization"],
            },
        });

        if (!response.ok) {
            throw new Error(
                "Network response was not ok: " + response.statusText,
            );
        }
        return response;
    }

    const addSourceBufferWhenOpen = (
        mediaSource: MediaSource,
        mimeStr: string,
    ) => {
        return new Promise<SourceBuffer>((res, rej) => {
            const getSourceBuffer = () => {
                try {
                    const sourceBuffer = mediaSource.addSourceBuffer(mimeStr);
                    res(sourceBuffer);
                } catch (e) {
                    rej(e);
                }
            };
            if (mediaSource.readyState === "open") {
                getSourceBuffer();
            } else {
                mediaSource.addEventListener("sourceopen", getSourceBuffer);
            }
        });
    };

    async function loadPreviewAudio() {
        if (loadState === "idle") {
            loadState = "loading";
            try {
                const streams = await GetStreams(track.urn);
                const res = await fetchAudioWithHeaders(
                    streams.Urls.preview_mp3_128_url,
                    streams.Headers,
                );

                const buf = await res.arrayBuffer();
                const mime = res.headers.get("content-type") ?? "";

                const mediaSource = new MediaSource();
                audioEl!.src = URL.createObjectURL(mediaSource);

                const sourceBuffer = await addSourceBufferWhenOpen(
                    mediaSource,
                    mime,
                );
                sourceBuffer.onupdateend = () => {
                    mediaSource.endOfStream();
                    audioEl!.play();
                };
                sourceBuffer.appendBuffer(buf);
                loadState = "done";
            } catch (e) {
                console.error(e);
                loadState = "error";
            }
        }
    }
</script>

<div class="flex flex-row items-center justify-center">
    {#if loadState === "idle" || loadState === "loading"}
        <Button size="icon" variant="ghost" onclick={loadPreviewAudio}>
            {#if loadState === "loading"}
                <Spinner />
            {/if}
            <Play />
        </Button>
    {:else}
        <Button
            size="icon"
            variant="ghost"
            onclick={() => {
                closed = !closed;
                if (closed) {
                    audioEl?.pause();
                } else {
                    audioEl?.play();
                }
            }}
        >
            {#if closed}
                <Play />
            {:else}
                <X />
            {/if}
        </Button>
    {/if}
    <audio
        bind:this={audioEl}
        volume={audioLevel}
        controls
        class={loadState === "done"
            ? closed
                ? "hidden"
                : "visible mx-auto pt-2"
            : "hidden"}
    ></audio>
</div>

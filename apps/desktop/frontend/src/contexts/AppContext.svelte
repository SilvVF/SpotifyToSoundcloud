<script lang="ts">
    import {
        QueryClient,
        QueryClientProvider,
    } from "@sveltestack/svelte-query";
    import { onDestroy, setContext } from "svelte";
    import { navigate } from "svelte5-router";
    import { SoundCloudAuthUrl, SpotifyAuthUrl } from "wailsjs/go/main/App";
    import SpotifyIcon from "$lib/icons/SpotifyIcon.svelte";
    import SoundcloudIcon from "$lib/icons/SoundcloudIcon.svelte";
    import { toast, Toaster } from "svelte-sonner";
    import type { AppContext, AuthEvent } from "src/types/types.svelte";
    import { EventsOn } from "wailsjs/runtime/runtime";

    const { children } = $props();

    const queryClient = new QueryClient({
        defaultOptions: {
            queries: {
                refetchOnWindowFocus: false,
            },
        },
    });

    const appState = $state<AppContext>({
        spotify: {
            name: "spotify",
            authed: false,
            err: undefined,
            authUrl: "",
            loading: true,
            icon: SpotifyIcon,
        },
        soundcloud: {
            name: "soundcloud",
            authed: false,
            err: undefined,
            authUrl: "",
            loading: true,
            icon: SoundcloudIcon,
        },
    });

    const unsub = EventsOn("auth_event", (data: AuthEvent) => {
        try {
            appState[data.name].authed = data.ok;
            appState[data.name].err = data.ok ? undefined : data.err;
        } catch (e) {
            toast.error(`failed to authorize ${data.name} err: ${data.err}`);
        }
    });

    onDestroy(() => unsub());

    $effect(() => {
        SpotifyAuthUrl().then((authUrl) => {
            appState.spotify.authUrl = authUrl;
            appState.spotify.loading = false;
        });
        SoundCloudAuthUrl().then((authUrl) => {
            appState.soundcloud.authUrl = authUrl;
            appState.soundcloud.loading = false;
        });
    });

    $effect(() => {
        if (appState.soundcloud.authed && appState.spotify.authed) {
            navigate("/home");
        }
    });

    setContext("appContextName", appState);
</script>

<QueryClientProvider client={queryClient}>
    <Toaster />
    {@render children?.()}
</QueryClientProvider>

<script lang="ts">
    import { Router, Route, navigate } from "svelte5-router";

    import { onMount, setContext } from "svelte";
    import { type AppState, type AuthEvent } from "./State";
    import { EventsOn, LogDebug } from "wailsjs/runtime";
    import {
        SoundCloudAuthUrl,
        SpotifyAuthUrl,
        RefreshAuthState,
    } from "wailsjs/go/main/App";
    import SpotifyIcon from "$lib/icons/SpotifyIcon.svelte";
    import SoundcloudIcon from "$lib/icons/SoundcloudIcon.svelte";
    import SignIn from "./routes/SignIn.svelte";
    import Home from "./routes/Home.svelte";
    import {
        QueryClient,
        QueryClientProvider,
    } from "@sveltestack/svelte-query";

    const queryClient = new QueryClient();

    const appState = $state<AppState>({
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

    onMount(() => {
        const unsub = EventsOn("auth_event", (data) => {
            LogDebug(`received auth event ${data}`);
            try {
                const value: AuthEvent = JSON.parse(data);
                appState[value.name].authed = value.ok;
                appState[value.name].err = value.ok ? undefined : value.err;
            } catch (e) {
                LogDebug(`authEvent: ${e}`);
            }
        });

        return () => unsub();
    });

    onMount(() => {
        SpotifyAuthUrl().then((authUrl) => {
            appState.spotify.authUrl = authUrl;
            appState.spotify.loading = false;
        });
        SoundCloudAuthUrl().then((authUrl) => {
            appState.soundcloud.authUrl = authUrl;
            appState.soundcloud.loading = false;
        });
        RefreshAuthState().catch((e) => LogDebug(e));
    });

    setContext("app_state", appState);
</script>

<QueryClientProvider client={queryClient}>
    <main>
        {#if appState.soundcloud.authed && appState.spotify.authed}
            <Home></Home>
        {:else}
            <SignIn />
        {/if}
    </main>
</QueryClientProvider>

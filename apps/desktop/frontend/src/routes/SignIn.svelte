<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import LoginCheck from "$lib/LoginCheck.svelte";
    import type { AppState } from "src/State";
    import { getContext } from "svelte";
    import { navigate } from "svelte5-router";

    const appState = getContext<AppState>("app_state");

    const ableToContinue = $derived(
        appState.soundcloud.authed && appState.spotify.authed,
    );

    function handleContinueClick() {
        if (ableToContinue) {
            navigate("/home");
        }
    }
</script>

<div class="text-2xl font-semibold">Sign into required platforms</div>
<div class="card flex flex-col space-y-4 w-full items-center">
    <LoginCheck state={appState.spotify} />
    <LoginCheck state={appState.soundcloud} />
    <Button onclick={handleContinueClick} disabled={!ableToContinue}
        >Continue</Button
    >
</div>

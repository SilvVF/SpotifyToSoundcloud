<script lang="ts">
    import { Button } from "$lib/components/ui/button";
    import LoginCheck from "$lib/LoginCheck.svelte";
    import { getContext } from "svelte";
    import { navigate } from "svelte5-router";
    import { RefreshAuthState } from "wailsjs/go/main/App";
    import { toast } from "svelte-sonner";
    import type { AppContext } from "src/types/types.svelte";

    const appState = getContext<AppContext>("appContextName");

    const ableToContinue = $derived(
        appState.soundcloud.authed && appState.spotify.authed,
    );

    function handleContinueClick() {
        if (ableToContinue) {
            navigate("/home");
        }
    }

    function handleRefresh() {
        RefreshAuthState().catch((e) => {
            toast.error(`failed to load auth state err: ${e}`);
        });
    }
</script>

<div class="flex flex-col items-center justify-center h-full w-full space-y-6">
    <div class="text-2xl font-semibold">Sign into required platforms</div>
    <div class="card flex flex-col space-y-4 w-full items-center">
        <LoginCheck state={appState.spotify} />
        <LoginCheck state={appState.soundcloud} />
        <Button onclick={handleContinueClick} disabled={!ableToContinue}>
            Continue
        </Button>
        {#if ableToContinue}
            <Button onclick={handleRefresh} variant="outline">Resfresh</Button>
        {/if}
    </div>
</div>

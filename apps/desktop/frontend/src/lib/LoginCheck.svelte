<script lang="ts">
    import type { ApiState } from "../State";
    import * as Item from "$lib/components/ui/item/index.js";
    import Button from "./components/ui/button/button.svelte";
    import { BrowserOpenURL } from "wailsjs/runtime/runtime";

    const { state }: { state: ApiState } = $props();

    function handleClick() {
        if (!state.authed && state.authUrl !== "") {
            BrowserOpenURL(state.authUrl);
        }
    }
</script>

<div class="flex flex-col">
    <Button
        class="w-xs"
        variant={state.authed ? "default" : "outline"}
        onclick={handleClick}
    >
        <Item.Root class="w-full">
            {#if !state.authed}
                <Item.Content>
                    <Item.Title class="line-clamp-1">
                        Sign in with {state.name}
                    </Item.Title>
                </Item.Content>
            {:else}
                <Item.Content>
                    <Item.Title class="line-clamp-1">
                        Signed into {state.name}
                    </Item.Title>
                </Item.Content>
            {/if}
            <Item.Content class="flex-none justify-end">
                <state.icon />
            </Item.Content>
        </Item.Root>
    </Button>
</div>

<script lang="ts">
    import { ModeWatcher } from "mode-watcher";
    import { toggleMode } from "mode-watcher";
    import { setContext, type Snippet } from "svelte";
    import type { ThemeContext } from "src/types/types.svelte";

    const { children }: { children: Snippet } = $props();

    let theme = $state<ThemeContext>({
        toggle: toggleMode,
    });

    setContext("themeContextName", theme);
</script>

<div>
    <ModeWatcher />
    {@render children?.()}
</div>

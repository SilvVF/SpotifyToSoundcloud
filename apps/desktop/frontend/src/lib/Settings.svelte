<script lang="ts">
    import { Button } from "$lib/components/ui/button/index.js";
    import * as Dialog from "$lib/components/ui/dialog/index.js";
    import { useLocation } from "svelte5-router";
    const location = useLocation();

    const settingsOpen = $derived.by(() => {
        if ($location?.state) {
            return Object.keys($location.state).includes("settings_open");
        }
        return false;
    });

    function handleOpenChange(open: boolean | undefined) {
        if (!open) {
            window.history.back();
        }
    }
</script>

<Dialog.Root open={settingsOpen} onOpenChange={handleOpenChange}>
    <Dialog.Content class="max-w-md h-[80%] overflow-auto">
        <Dialog.Header>
            <Dialog.Title>Settings</Dialog.Title>
            <Dialog.Description>
                Make changes to your settings here. Click save when you're done.
            </Dialog.Description>
        </Dialog.Header>
        <div class="w-full h-full"></div>
        <Dialog.Footer>
            <Button type="submit">Save changes</Button>
        </Dialog.Footer>
    </Dialog.Content>
</Dialog.Root>

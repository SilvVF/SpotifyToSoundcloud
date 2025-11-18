<script lang="ts">
    import * as Menubar from "$lib/components/ui/menubar/index.js";
    import { setMode, mode, type SystemModeValue } from "mode-watcher";
    import { navigate } from "svelte5-router";


    let bookmarks = $state(false);
    let fullUrls = $state(true);

    function handleSettingsSelected() {
        const state = {
            "settings_open": "true",
        };
        navigate(window.location.toString(), { replace: false, state });
    }

    function handleGoBack() {
        window.history.back();
    }

    function handleGoHome() {
        navigate("/home")
    }
</script>

<Menubar.Root>
    <Menubar.Menu>
        <Menubar.Trigger>File</Menubar.Trigger>
        <Menubar.Content>
            <Menubar.Item onclick={handleGoBack}>Back</Menubar.Item>
            <Menubar.Item onclick={handleGoHome}>Home</Menubar.Item>
            <Menubar.Item onclick={handleSettingsSelected}>
                Settings <Menubar.Shortcut>⌘Alt+S</Menubar.Shortcut>
            </Menubar.Item>
            <Menubar.Separator />
            <Menubar.Sub>
                <Menubar.SubTrigger>Share</Menubar.SubTrigger>
                <Menubar.SubContent>
                    <Menubar.Item>Email link</Menubar.Item>
                    <Menubar.Item>Messages</Menubar.Item>
                    <Menubar.Item>Notes</Menubar.Item>
                </Menubar.SubContent>
            </Menubar.Sub>
            <Menubar.Separator />
            <Menubar.Item>
                Print... <Menubar.Shortcut>⌘P</Menubar.Shortcut>
            </Menubar.Item>
        </Menubar.Content>
    </Menubar.Menu>
    <Menubar.Menu>
        <Menubar.Trigger>View</Menubar.Trigger>
        <Menubar.Content>
            <Menubar.CheckboxItem bind:checked={bookmarks}
                >Always Show Bookmarks Bar</Menubar.CheckboxItem
            >
            <Menubar.CheckboxItem bind:checked={fullUrls}>
                Always Show Full URLs
            </Menubar.CheckboxItem>
            <Menubar.Separator />
            <Menubar.Item inset>
                Reload <Menubar.Shortcut>⌘R</Menubar.Shortcut>
            </Menubar.Item>
            <Menubar.Item inset>
                Force Reload <Menubar.Shortcut>⇧⌘R</Menubar.Shortcut>
            </Menubar.Item>
            <Menubar.Separator />
            <Menubar.Item inset>Toggle Fullscreen</Menubar.Item>
            <Menubar.Separator />
            <Menubar.Item inset>Hide Sidebar</Menubar.Item>
        </Menubar.Content>
    </Menubar.Menu>
    <Menubar.Menu>
        <Menubar.Trigger>Theme</Menubar.Trigger>
        <Menubar.Content>
            <Menubar.RadioGroup
                value={mode.current}
                onValueChange={(item) =>
                    setMode((item as SystemModeValue) ?? "system")}
            >
                <Menubar.RadioItem value="light">Light</Menubar.RadioItem>
                <Menubar.RadioItem value="dark">Dark</Menubar.RadioItem>
                <Menubar.RadioItem value="">System</Menubar.RadioItem>
            </Menubar.RadioGroup>
        </Menubar.Content>
    </Menubar.Menu>
</Menubar.Root>

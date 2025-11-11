<script lang="ts">
    import Button from "$lib/components/ui/button/button.svelte";
    import { Link } from "svelte5-router";

    const { playlistId }: { playlistId: string } = $props();
</script>

<div class="flex flex-col h-full w-full scroll-auto">
    <Link to="/home">
        <Button>Back</Button>
    </Link>
    <div>Convert {playlistId}</div>
</div>

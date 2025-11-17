<script lang="ts">
    import { Router, Route } from "svelte5-router";
    import SignIn from "./routes/SignIn.svelte";
    import Home from "./routes/Home.svelte";
    import Convert from "./routes/Convert.svelte";
    import AppContext from "./contexts/AppContext.svelte";
    import ThemeContext from "./contexts/ThemeContext.svelte";
    import GenerateContext from "./contexts/GenerateContext.svelte";
    import AppMenubar from "$lib/AppMenubar.svelte";
    import Settings from "$lib/Settings.svelte";

    let url = $state("/signin");
</script>

<main>
    <AppContext>
        <GenerateContext>
            <ThemeContext>
                <div class="w-screen h-screen flex flex-col overflow-clip">
                    <AppMenubar />
                    <Router {url}>
                        <Settings />
                        <Route path="/home" component={Home} />
                        <Route path="/signin" component={SignIn} />
                        <Route path="/convert/:id">
                            {#snippet children(params)}
                                {#key params.id}
                                    <Convert playlistId={params.id} />
                                {/key}
                            {/snippet}
                        </Route>
                    </Router>
                </div>
            </ThemeContext>
        </GenerateContext>
    </AppContext>
</main>

<script>
    import {Router, Link, Route} from "svelte-routing";
    import Dashboard from "./routes/Dashboard.svelte";
    import Header from "./layout/Header.svelte";
    import Instances from "./routes/Instances.svelte";
    import Config from "./routes/Config.svelte";
    import SiteDetails from "./routes/SiteDetails.svelte";

    export let url = "";
</script>

<Router url="{url}">
    <Header />
    <div class="container-xl">
        <Route path="/"><Dashboard /></Route>
        <Route path="/config"><Config /></Route>
        <Route path="/instances"><Instances /></Route>
        <Route path="/sites/:name" let:params><SiteDetails name="{params.name}" /></Route>
    </div>
</Router>
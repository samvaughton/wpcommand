<script>
    import {Router, Link, Route} from "svelte-routing";
    import Dashboard from "./routes/Dashboard.svelte";
    import Header from "./layout/Header.svelte";
    import Instances from "./routes/Instances.svelte";
    import Config from "./routes/Config.svelte";
    import SiteDetails from "./routes/SiteDetails.svelte";
    import Login from "./routes/Login.svelte";
    import { userStore } from './store/user'

    export let url = "";

    if (!$userStore) {
        url = "/login"
    }
</script>

<Router url="{url}">
    <Header />
    <div class="container-xl">
        <Route path="/login"><Login /></Route>
        <Route path="/"><Dashboard /></Route>
        <Route path="/config"><Config /></Route>
        <Route path="/instances"><Instances /></Route>
        <Route path="/sites/:key" let:params><SiteDetails key="{params.key}" /></Route>
    </div>
</Router>